package sunspec

import (
	"context"

	"github.com/GoAethereal/modbus"
)

// NewServer creates a new sunspec server with the given configuration.
func NewServer(cfg Config) *Server {
	return &Server{server: newModbusServer(cfg.Endpoint, cfg.logger()), logger: cfg.logger()}
}

// Server is a sunspec compliant server.
type Server struct {
	server
	models Models
	logger Logger
}

var _ Device = (*Server)(nil)

// Model returns the first model identifies by id.
func (s *Server) Model(id uint16) Model { return s.models[1 : len(s.models)-1].Model(id) }

// Models returns all models from the device.
func (s *Server) Models(ids ...uint16) Models { return s.models[1 : len(s.models)-1].Models(ids...) }

// Serve instantiates the model, as declared in the definition and starts serving it to connected clients.
// The handler function is called for any incoming client request.
func (s *Server) Serve(ctx context.Context, handler func(ctx context.Context, isWrite bool, pts Points) error, defs ...Definition) error {
	// append the start marker
	s.models = append(Models(nil), marker(0))
	adr := ceil(s.models.First())
	for _, def := range defs {
		s.logger.Info("instantiating model definition", def.ID(), "at address", adr)
		m, err := def.Instance(adr, func(pts []Point) error { return nil })
		if err != nil {
			return err
		}
		s.logger.Info("verifying model", def.ID())
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
	serve(ctx context.Context, d Device, handler func(ctx context.Context, isWrite bool, pts Points) error) error
}

var _ server = (*mbServer)(nil)

type mbServer struct {
	mb     modbus.Server
	opt    modbus.Options
	logger Logger
}

func newModbusServer(endpoint string, l Logger) *mbServer {
	return &mbServer{
		mb: modbus.Server{},
		opt: modbus.Options{
			Mode:     "tcp",
			Kind:     "tcp",
			Endpoint: endpoint,
		},
		logger: l,
	}
}

func (s *mbServer) serve(ctx context.Context, d Device, handler func(ctx context.Context, isWrite bool, pts Points) error) error {
	return s.mb.Serve(ctx, s.opt, &modbus.Mux{
		ReadHoldingRegisters: func(ctx context.Context, address, quantity uint16) (res []byte, ex modbus.Exception) {
			s.logger.Debug("received modbus read request for address", address, "with quantity", quantity)
			pts, err := collect(d, index{address: address, quantity: quantity})
			if err != nil {
				return nil, modbus.ExIllegalDataAddress
			}
			if err := handler(ctx, false, pts); err != nil {
				return nil, modbus.ExSlaveDeviceFailure
			}
			buf := make([]byte, 2*pts.Quantity())
			if err := pts.encode(buf); err != nil {
				return nil, modbus.ExSlaveDeviceFailure
			}
			return buf, nil
		},
		WriteMultipleRegisters: func(ctx context.Context, address uint16, values []byte) (ex modbus.Exception) {
			s.logger.Debug("received modbus write request for address", address, "with payload", values)
			pts, err := collect(d, index{address: address, quantity: uint16(len(values) * 2)})
			if err != nil {
				return modbus.ExIllegalDataAddress
			}
			// ref 6.5.1 / 6.5.3: Unimplemented Registers / Writing a Read-Only Register
			for _, p := range pts {
				if !p.Valid() || !p.Writable() {
					return modbus.ExIllegalDataAddress
				}
			}
			if err := pts.decode(values); err != nil {
				return modbus.ExSlaveDeviceFailure
			}
			if err := handler(ctx, true, pts); err != nil {
				return modbus.ExSlaveDeviceFailure
			}
			return nil
		},
	})
}
