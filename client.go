package sunspec

import (
	"errors"

	"github.com/GoAethereal/cancel"
	"github.com/GoAethereal/modbus"
)

// Client represents a compliant sunspec client.
type Client struct {
	client
	Device
}

var _ Device = (*Client)(nil)

// Scan analyses the server retrieving its device.
// The process uses the given definition as reference.
func (c *Client) Scan(ctx cancel.Context, defs ...Definition) (err error) {
	c.Device, err = c.scan(ctx, defs)
	return err
}

// Read requests all point values in the given address range from the server.
func (c *Client) Read(ctx cancel.Context, idx ...Index) (Points, error) {
	pts, err := collect(c, idx...)
	if err != nil {
		return nil, err
	}
	return c.read(ctx, pts...)
}

// Write sends all point values in the given address range to the server.
// Read-Only points are silently skipped.
func (c *Client) Write(ctx cancel.Context, idx ...Index) (Points, error) {
	pts, err := collect(c, idx...)
	if err != nil {
		return nil, err
	}
	// filter out writable points ignore read only
	var i int
	for _, p := range pts {
		if p.Writable() {
			pts[i] = p
			i++
		}
	}
	pts = pts[:i]
	if len(pts) == 0 {
		return nil, errors.New("sunspec: no writable points for given index")
	}
	return c.write(ctx, pts...)
}

type client interface {
	// Connect starts the underlying server-connection.
	Ready() bool
	// Disconnect stops the underlying server-connection.
	Disconnect()
	scan(ctx cancel.Context, defs []Definition) (Device, error)
	read(ctx cancel.Context, pts ...Point) (Points, error)
	write(ctx cancel.Context, pts ...Point) (Points, error)
}

var _ client = (*mbClient)(nil)

type mbClient struct {
	modbus.Client
}

func newModbusClient(endpoint string) *mbClient {
	return &mbClient{
		Client: modbus.Client{Config: modbus.Config{
			Mode:     "tcp",
			Kind:     "tcp",
			Endpoint: endpoint,
		}},
	}
}

func (c *mbClient) scan(ctx cancel.Context, defs []Definition) (Device, error) {
	adr, err := c.marker(ctx)
	if err != nil {
		return nil, err
	}
	var (
		m Model
		d Models
	)
	h := header(adr+2, 0, 0)
	for {
		if _, err := c.read(ctx, h.Points()...); err != nil {
			return nil, err
		}
		if h.ID().Get() == 0xFFFF {
			return d, nil
		}
		m = h
		for _, def := range defs {
			if def.ID() == h.ID().Get() {
				m, err = def.Instance(h.Address(), func(pts []Point) error {
					_, err := c.read(ctx, pts...)
					return err
				})
				if err != nil {
					return nil, err
				}
				if err := Verify(m); err != nil {
					return nil, err
				}
				break
			}
		}
		d = append(d, m)
		h = header(h.Address()+h.Length().Get()+2, 0, 0)
	}
}

// marker locates the modbus stating address of the endpoint by scanning the base addresses.
func (c *mbClient) marker(ctx cancel.Context) (uint16, error) {
	for _, adr := range [...]uint16{0, 40000, 50000} {
		if _, err := c.read(ctx, marker(adr).Points()...); err == nil {
			return adr, nil
		}
	}
	return 0, errors.New("sunspec: could not identify the starting marker")
}

// read attempts to request the data for all given points from the modbus endpoint.
func (c *mbClient) read(ctx cancel.Context, pts ...Point) (Points, error) {
	return c.execute(125, pts, func(pts Points) error {
		res, err := c.ReadHoldingRegisters(ctx, pts.address(), pts.Quantity())
		if err != nil {
			return err
		}
		return pts.decode(res)
	})
}

// write attempts to send the point values of all given points to the modbus endpoint
func (c *mbClient) write(ctx cancel.Context, pts ...Point) (Points, error) {
	return c.execute(123, pts, func(pts Points) error {
		req := make([]byte, 2*pts.Quantity())
		if err := pts.encode(req); err != nil {
			return err
		}
		return c.WriteMultipleRegisters(ctx, pts.address(), req)
	})
}

// execute calls back cmd for all given points.
// The input collection is split in regards to their modbus continuity limited by the given register limit.
// 	ToDo: still needs handling for sync groups
func (c *mbClient) execute(limit uint16, pts Points, cmd func(pts Points) error) (Points, error) {
	for i, j, l := 1, 0, len(pts); j < l; j = i {
		for _, p := range pts[i:] {
			if ceil(pts[i-1]) != p.Address() || ceil(p)-pts[j].Address() > limit {
				break
			}
			i++
		}
		if err := cmd(pts[j:i]); err != nil {
			return pts[:j], err
		}
	}
	return pts, nil
}
