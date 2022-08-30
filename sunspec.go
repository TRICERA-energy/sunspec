package sunspec

// marker returns a dummy model for representing the magic identifier SunS.
func marker(adr uint16) Model {
	return &model{
		&group{
			name: "marker",
			points: Points{
				&tString{
					data: []byte("SunS"),
					point: point{
						name:    "SunS",
						static:  true,
						address: adr,
					},
				},
			},
		},
	}
}

// header returns a prototype for identifying a model using the minimum requirements.
func header(adr, id, l uint16) Model {
	return &model{
		&group{
			name: "header",
			points: Points{
				&tUint16{
					data: id,
					point: point{
						name:    "ID",
						static:  true,
						address: adr,
					},
				},
				&tUint16{
					data: l,
					point: point{
						name:    "L",
						static:  true,
						address: adr + 1,
					},
				},
			},
		},
	}
}
