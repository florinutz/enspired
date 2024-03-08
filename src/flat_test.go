package src

import (
	"strings"
	"testing"
)

func TestRoomParser_Ingest(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *FlatParser
		wantErr bool
	}{
		{
			name:    "Upper perimeter",
			input:   "+-----------+------------------------------------+",
			want:    &FlatParser{Line: 1},
			wantErr: false,
		},
		{
			name: "2 open rooms",
			input: `
+----+---+
|    |   |`,
			want: &FlatParser{
				Line: 3,
				OpenRooms: []*openRoom{
					{
						RoomData: newRoomData(),
						segments: NewLineSegments(&segment{
							start:   1,
							content: "      ",
						}),
					},
				},
				closedRooms: []*roomData{},
			},
			wantErr: false,
		},
		{
			name: "1 closed room",
			input: `
+------+
| W    |
+------+
`,
			want: &FlatParser{
				Line: 5,
				closedRooms: []*roomData{
					{
						Chairs: map[rune]int{'W': 1},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "2 rooms, titles",
			input: `
+------+
| W    |-----+
+------+ P   |
|   C        |
|     P      |
|   ( room)  |
+------------+
`,
			want: &FlatParser{
				Line: 9,
				closedRooms: []*roomData{
					{
						Chairs: map[rune]int{'W': 1},
					},
					{
						Name:   "room",
						Chairs: map[rune]int{'P': 2, 'C': 2},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one room inside the other",
			input: `
+----------------+
| W          C   |
| C +------+     |
+---+W (x2)|     |
    +-+----+   P |
      |          |
      |  +-------+
+-----+  |
| P      +-----+
|              |
+-----+-----+  |
            |  |
  +---------+  |
  | (room) W W |
  +------------+
`,
			want: &FlatParser{
				Line: 17,
				closedRooms: []*roomData{
					{
						Name:   "x2",
						Chairs: map[rune]int{'W': 1},
					},
					{
						Name:   "room",
						Chairs: map[rune]int{'W': 1, 'C': 2, 'P': 1},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		parser := NewRoomParser()
		t.Run(tt.name, func(t *testing.T) {
			// Split the lines by newline character
			individualLines := strings.Split(tt.input, "\n")

			// Ingest each line
			for _, line := range individualLines {
				if err := parser.Ingest(line); (err != nil) != tt.wantErr {
					t.Fatalf("Parser: error:\n\t%+v\nwant err: %+v\n\ninput:\n%s",
						err, tt.wantErr, tt.input)
				}
			}

			if parser.Line != tt.want.Line {
				t.Fatalf("Parser: line = %d, want %d", parser.Line, tt.want.Line)
			}

			if len(parser.closedRooms) != len(tt.want.closedRooms) {
				t.Fatalf("Parser: number of closed rooms = %d, want %d",
					len(parser.closedRooms), len(tt.want.closedRooms))
			}

			for i, wantedClosedRoom := range tt.want.closedRooms {
				if wantedClosedRoom.Name != parser.closedRooms[i].Name {
					t.Errorf("Parser: room name = %s, want %s",
						parser.closedRooms[i].Name, wantedClosedRoom.Name)
				}
				for gotChairType, count := range parser.closedRooms[i].Chairs {
					if _, ok := parser.closedRooms[i].Chairs[gotChairType]; !ok {
						t.Fatalf("Parser: got extra chair type %c in closed room %d",
							gotChairType, i)
					}
					if count != parser.closedRooms[i].Chairs[gotChairType] {
						t.Fatalf("Parser: chair count = %d, want %d",
							count, tt.want.closedRooms[i].Chairs[gotChairType])
					}
				}
			}
		})
	}
}
func TestRoomString(t *testing.T) {
	tt := []struct {
		name string
		rd   *roomData
		want string
	}{
		{
			name: "Empty Case",
			rd:   &roomData{},
			want: "(no data)",
		},
		{
			name: "One Chair",
			rd: &roomData{
				Name:   "Living Room",
				Chairs: map[rune]int{'A': 1},
			},
			want: "Living Room:\nA: 1",
		},
		{
			name: "Multiple Chairs",
			rd: &roomData{
				Name:   "Bedroom",
				Chairs: map[rune]int{'A': 1, 'B': 2},
			},
			want: "Bedroom:\nA: 1, B: 2",
		},
		{
			name: "Multiple Chairs Unordered",
			rd: &roomData{
				Name:   "Kitchen",
				Chairs: map[rune]int{'B': 2, 'A': 1},
			},
			want: "Kitchen:\nA: 1, B: 2",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.rd.String()
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// the only benchmark
// todo: more (for Split maybe?)
func BenchmarkRoomParser_Ingest(b *testing.B) {
	parser := NewRoomParser()
	input := `
+----------------+
| W          C   |
| C +------+     |
+---+W (x2)|     |
    +-+----+   P |
      |          |
      |  +-------+
+-----+  |
| P      +-----+
|              |
+-----+-----+  |
            |  |
  +---------+  |
  | (room) W W |
+-+---------+--+---------------------------------+
|           |                                    |
| (closet)  |                                    |
|         P |                            S       |
|         P |         (sleeping room)            |
|         P |                                    |
|           |                                    |
+-----------+    W                               |
|           |                                    |
|        W  |                                    |
|           |                                    |
|           +--------------+---------------------+
|                          |                     |
|                          |                W W  |
|                          |    (office)         |
|                          |                     |
+--------------+           |                     |
|              |           |                     |
| (toilet)     |           |             P       |
|   C          |           |                     |
|              |           |                     |
+--------------+           +---------------------+
|              |           |                     |
|              |           |                     |
|              |           |                     |
| (bathroom)   |           |      (kitchen)      |
|              |           |                     |
|              |           |      W   W          |
|              |           |      W   W          |
|       P      +           |                     |
|             /            +---------------------+
|            /                                   |
|           /                                    |
|          /     +-----+              W    W   W |
+---------+      | (x) |                         |
|                +-----+                         |
| S                                   W    W   W |
|                (living room)                   |
| S                                              |
|                                                |
|                                                |
|                                                |
|                                                |
+--------------------------+---------------------+
                           |                     |
                           |                  P  |
                           |  (balcony)          |
                           |                 P   |
                           |                     |
                           +---------------------+
`
	individualLines := strings.Split(input, "\n")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, line := range individualLines {
			_ = parser.Ingest(line)
		}
	}
}
