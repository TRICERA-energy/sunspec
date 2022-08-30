package sunspec

import (
	"github.com/GoAethereal/cancel"
	"github.com/GoAethereal/modbus"
)

// Server is a sunspec compliant server.
type Server struct {
	server
	models Models
}

var _ Device = (*Server)(nil)

// Model returns the first model identifies by id.
func (s *Server) Model(id uint16) Model { return s.models[1 : len(s.models)-1].Model(id) }

// Models returns all models from the device.
func (s *Server) Models(ids ...uint16) Models { return s.models[1 : len(s.models)-1].Models(ids...) }

// Serve instantiates the model, as declared in the definition and starts serving it to connected clients.
// The handler function is called for any incoming client request.
func (s *Server) Serve(ctx cancel.Context, handler func(ctx cancel.Context, req Request) error, defs ...Definition) error {
	// append the start marker
	s.models = append(Models(nil), marker(0))
	adr := ceil(s.models.First())
	for _, def := range defs {
		m, err := def.Instance(adr, func(pts []Point) error { return nil })
		if err != nil {
			return err
		}
		if err := Verify(m); err != nil {
			return err
		}
		adr = ceil(m)
		s.models = append(s.models, m)
	}
	// append the endmarker
	s.models = append(s.models, header(adr, 0xFFFF, 0))

	return s.serve(ctx, s.models, handler)
}

type server interface {
	serve(ctx cancel.Context, d Device, handler func(ctx cancel.Context, req Request) error) error
}

var _ server = (*mbServer)(nil)

type mbServer struct {
	modbus.Server
}

func newModbusServer(endpoint string) *mbServer {
	return &mbServer{
		Server: modbus.Server{Config: modbus.Config{
			Mode:     "tcp",
			Kind:     "tcp",
			Endpoint: endpoint,
		}},
	}
}

func (s *mbServer) serve(ctx cancel.Context, d Device, handler func(ctx cancel.Context, req Request) error) error {
	return s.Serve(ctx, &modbus.Mux{
		ReadHoldingRegisters: func(ctx cancel.Context, address, quantity uint16) (res []byte, ex modbus.Exception) {
			pts, err := collect(d, index{address: address, quantity: quantity})
			if err != nil {
				return nil, modbus.IllegalDataAddress
			}
			req := &request{points: pts, writing: false, buffer: make([]byte, 2*pts.Quantity())}
			if err := handler(ctx, req); err != nil {
				return nil, modbus.SlaveDeviceFailure
			}
			return req.buffer, 0
		},
		WriteMultipleRegisters: func(ctx cancel.Context, address uint16, values []byte) (ex modbus.Exception) {
			pts, err := collect(d, index{address: address, quantity: uint16(len(values) / 2)})
			if err != nil {
				return modbus.IllegalDataAddress
			}
			// ref 6.5.1 / 6.5.3: Unimplemented Registers / Writing a Read-Only Register
			for _, p := range pts {
				if !p.Valid() || !p.Writable() {
					return modbus.IllegalDataAddress
				}
			}
			req := &request{points: pts, writing: true, buffer: values}
			if err := handler(ctx, req); err != nil {
				return modbus.SlaveDeviceFailure
			}
			return 0
		},
	})
}
