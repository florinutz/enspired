package src

import (
	"reflect"
	"sort"
	"testing"
)

func TestSplitLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected LineSegments
	}{
		{
			name:     "Empty Line",
			line:     "",
			expected: NewLineSegments(),
		},
		{
			name:     "No Delimiters",
			line:     "HelloWorld",
			expected: NewLineSegments(),
		},
		{
			name:     "Single Delimiter (still no room)",
			line:     "Hello|World",
			expected: NewLineSegments(),
		},
		{
			name:     "Multiple Delimiters, one room",
			line:     "Hello/World|Rooms",
			expected: NewLineSegments(&segment{6, "World"}),
		},
		{
			name: "2 rooms",
			line: "|  |   |",
			expected: NewLineSegments(
				&segment{1, "  "},
				&segment{4, "   "},
			),
		},
		{
			name: "2 rooms, extra stuff along",
			line: "  |  |   |   ",
			expected: NewLineSegments(
				&segment{3, "  "},
				&segment{6, "   "},
			),
		},
		{
			name:     "balcony line from example",
			line:     "                           |  (balcony)          |",
			expected: NewLineSegments(&segment{28, "  (balcony)          "}),
		},
		{
			name:     "line 12 from example",
			line:     "|           +--------------+---------------------+",
			expected: NewLineSegments(&segment{1, "           "}),
		},
		{
			name: "line 17 from example",
			line: "+--------------+           |                     |",
			expected: NewLineSegments(
				&segment{16, "           "},
				&segment{28, "                     "},
			),
		},
		{
			name:     "line 22 from example",
			line:     "+--------------+           +---------------------+",
			expected: NewLineSegments(&segment{16, "           "}),
		},
		{
			name: "line 31 from example",
			line: "|             /            +---------------------+",
			expected: NewLineSegments(
				&segment{1, "             "},
				&segment{15, "            "},
			),
		},
		{
			name: "surrounding",
			line: "| +-+ |",
			expected: NewLineSegments(
				&segment{1, " "},
				&segment{5, " "},
			),
		},
	}

	// Loop through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function with the input provided by test case
			got := Split(tc.line)

			// Check if the result matches the expected result
			if !compareSegmentsLists(got, tc.expected) {
				t.Errorf("splitLine failed for %v, expected %v, got %v", tc.name, tc.expected, got)
			}
		})
	}
}

func compareSegmentsLists(s1, s2 LineSegments) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i].start != s2[i].start || s1[i].content != s2[i].content {
			return false
		}
	}
	return true
}

func TestOverlaps(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		length   int
		segments LineSegments
		want     LineSegments
	}{
		{
			name:     "Empty Segments",
			start:    3,
			length:   10,
			segments: NewLineSegments(),
			want:     NewLineSegments(),
		},
		{
			name:   "NoOverlap",
			start:  0,
			length: 5,
			segments: NewLineSegments(
				&segment{6, "Segment"},
			),
			want: NewLineSegments(),
		},
		{
			name:   "Overlap Start",
			start:  2,
			length: 3,
			segments: NewLineSegments(
				&segment{0, "bae"},
			),
			want: NewLineSegments(
				&segment{0, "bae"},
			),
		},
		{
			name:   "Overlap End",
			start:  0,
			length: 3,
			segments: NewLineSegments(
				&segment{2, "bae"},
			),
			want: NewLineSegments(
				&segment{2, "bae"},
			),
		},
		{
			name:   "Short Segment",
			start:  3,
			length: 1,
			segments: NewLineSegments(
				&segment{0, "abc"},
				&segment{3, "def"},
			),
			want: NewLineSegments(
				&segment{3, "def"},
			),
		},
		{
			name:   "SurroundedSegment",
			start:  2,
			length: 2,
			segments: NewLineSegments(
				&segment{0, "abc"},
				&segment{3, "def"},
				&segment{6, "ghi"},
			),
			want: NewLineSegments(
				&segment{0, "abc"},
				&segment{3, "def"},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Overlaps(tt.start, tt.length, tt.segments); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Overlaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultipleOverlaps(t *testing.T) {
	// specifies a set of segments from a line
	set := func(line string, indexes ...int) (ret LineSegments) {
		if line == "" {
			return nil
		}
		spl := Split(line)
		if len(indexes) == 0 {
			return spl
		}
		for i, seg := range spl {
			for _, index := range indexes {
				if i == index {
					ret = append(ret, seg)
				}
			}
		}
		return
	}

	tests := []struct {
		name               string
		set1               LineSegments
		set2               LineSegments
		wantOverlapping    LineSegments
		wantNonOverlapping LineSegments
	}{
		{
			name:               "Test1",
			set1:               set("|   |          |    |", 0, 2),
			set2:               set("|     |  |  |      |", 1, 2),
			wantOverlapping:    nil,
			wantNonOverlapping: set("|     |  |  |      |", 1, 2),
		},
		{
			name:               "Overlap",
			set1:               set("|   |          |    |", 2),
			set2:               set("|     |  |  |      |", 0, 1, 2, 3),
			wantOverlapping:    set("|     |  |  |      |", 3),
			wantNonOverlapping: set("|     |  |  |      |", 0, 1, 2),
		},
		{
			name:               "Bug",
			set1:               set("| C +------+     |", 0, 1),
			set2:               set("+---+ W    |     |", 1),
			wantOverlapping:    set("+---+ W    |     |", 1),
			wantNonOverlapping: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOverlapping, gotNonOverlapping := MultipleOverlaps(tt.set1, tt.set2)
			if !compareSegmentSets(gotOverlapping, tt.wantOverlapping) {
				t.Fatalf("MultipleOverlaps() overlapping segments = %+v, want %+v", gotOverlapping, tt.wantOverlapping)
			}
			if !compareSegmentSets(gotNonOverlapping, tt.wantNonOverlapping) {
				t.Fatalf("MultipleOverlaps() nonOverlappingSegments = %+v, want %+v", gotNonOverlapping, tt.wantNonOverlapping)
			}
		})
	}
}

func compareSegmentSets(set1, set2 LineSegments) bool {
	if len(set1) != len(set2) {
		return false
	}
	sort.Slice(set1, func(i, j int) bool { return set1[i].start < set1[j].start })
	sort.Slice(set2, func(i, j int) bool { return set2[i].start < set2[j].start })
	for i := range set1 {
		if set1[i].start != set2[i].start || set1[i].content != set2[i].content {
			return false
		}
	}
	return true
}

func TestSegmentData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantRoom *roomData
		wantErr  bool
	}{
		{
			name:     "room with all elements",
			input:    "(Living Room) WPSSC",
			wantRoom: &roomData{Name: "Living Room", Chairs: map[rune]int{'W': 1, 'P': 1, 'S': 2, 'C': 1}},
			wantErr:  false,
		},
		{
			name:    "title on multiple lines",
			input:   "(Living room WPSC",
			wantErr: true,
		},
		{
			// X is not a chair and neither is Z, but Z is ok coz it's in the title
			name:    "invalid character",
			input:   "W(PZ)C X",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoom, err := segmentData(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("segmentData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if gotRoom.Name != tt.wantRoom.Name {
				t.Errorf("got room.Name = %v, want %v", gotRoom.Name, tt.wantRoom.Name)
			}

			for chair, qty := range tt.wantRoom.Chairs {
				if gotRoom.Chairs[chair] != qty {
					t.Errorf("got room.Chairs[%v] = %v, want %v", string(chair),
						gotRoom.Chairs[chair], qty)
				}
			}
		})
	}
}
