package src

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// the data that we're actually interested in
type roomData struct {
	Name   string
	Chairs map[rune]int
}

// The string representation of a room's data.
// It shows the chairs ordered alphabetically.
//
// living room:
// W: 3, P: 0, S: 0, C: 0
func (d *roomData) String() string {
	if len(d.Chairs) == 0 {
		if d.Name == "" {
			return "(no data)"
		}
		return d.Name
	}
	var pairs []string
	keys := make([]rune, 0, len(d.Chairs))
	for k := range d.Chairs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		pairs = append(pairs, fmt.Sprintf("%c: %d", key, d.Chairs[key]))
	}
	return fmt.Sprintf("%s:\n%s", d.Name, strings.Join(pairs, ", "))
}

func (d *roomData) appendDataFromSegments(segments LineSegments) error {
	for _, segment := range segments {
		data, err := segmentData(segment.content)
		if err != nil {
			return fmt.Errorf("[segment: '%s'] error parsing segment: %w", segment.content, err)
		}
		d.append(data)
	}
	return nil
}

func (d *roomData) append(d2 *roomData) {
	if d2.Name != "" {
		d.Name = d2.Name
	}
	if d.Chairs == nil {
		d.Chairs = map[rune]int{}
	}
	for incomingChairType, incomingCount := range d2.Chairs {
		if _, exists := d.Chairs[incomingChairType]; !exists {
			d.Chairs[incomingChairType] = incomingCount
			continue
		}
		d.Chairs[incomingChairType] += incomingCount
	}
	return
}

// as the input and their segments keep coming, an open room will be one what was not yet closedRooms.
type openRoom struct {
	RoomData *roomData
	// RoomData segments as seen in the last parsed line.
	// These are needed for comparison with the segments of a new line.
	// The parser will take in line after line, and will record the evolution of the openRoom's set of segments,
	// so it will know at any time the exact state of an open room's segments
	segments LineSegments
}

type FlatParser struct {
	// the rooms that the parser currently knows as open.
	OpenRooms []*openRoom
	// rooms that were closed by walls of the perimeter or other rooms.
	// Closed rooms are what we're actually interested in.
	// All rooms will be closed by the end.
	closedRooms []*roomData
	Line        int
}

func NewRoomParser() *FlatParser {
	return &FlatParser{
		OpenRooms:   []*openRoom{},
		closedRooms: []*roomData{},
		Line:        0,
	}
}

func newRoomData() *roomData {
	return &roomData{Chairs: map[rune]int{}}
}

// Ingest parses input segments one by one and collects rooms data:
// - changes in the room walls positions (so we always know which segments belong to which open room)
// - room data (title, chairs, whatever interesting in the respective room segment)
func (p *FlatParser) Ingest(line string) error {
	p.Line++

	lineSegments := Split(line)

	for _, room := range p.OpenRooms {
		overlaps, rest := MultipleOverlaps(room.segments, lineSegments)
		if len(overlaps) == 0 {
			if err := p.closeRoom(room); err != nil {
				return fmt.Errorf("can't close room: %w", err)
			}
			continue
		}

		if err := room.RoomData.appendDataFromSegments(overlaps); err != nil {
			return fmt.Errorf("error ingesting segments: %w", err)
		}

		// keep the room's latest segments,
		// so we can compute overlaps in further loop iterations
		room.segments = overlaps

		// remove the room's segments from the current line's segments,
		// so they won't be parsed again in further loop iterations
		lineSegments = rest
	}

	// open rooms for each remaining (unassociated with previously opened rooms) lineSegment
	for _, segment := range lineSegments {
		data := &roomData{}
		if err := data.appendDataFromSegments(LineSegments{segment}); err != nil {
			return fmt.Errorf("[line %d] can't ingest segment: %w", p.Line, err)
		}
		p.OpenRooms = append(p.OpenRooms, &openRoom{
			RoomData: data,
			segments: LineSegments{segment},
		})
	}

	return nil
}

// IngestAllFromReader goes line by line, parses line segments, keeps the record of associations of current line segments with rooms
func (p *FlatParser) IngestAllFromReader(reader *bufio.Reader) error {
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	if err := p.Ingest(line); err != nil {
		return fmt.Errorf("[line %d] error parsing line: %w", p.Line, err)
	}

	return p.IngestAllFromReader(reader)
}

func (p *FlatParser) closeRoom(room *openRoom) error {
	for i, r := range p.OpenRooms {
		if r != room {
			continue
		}
		p.closedRooms = append(p.closedRooms, p.OpenRooms[i].RoomData)
		p.OpenRooms = append(p.OpenRooms[:i], p.OpenRooms[i+1:]...)
		break
	}

	return nil
}

func (p *FlatParser) totals(totalsEntryName string) *roomData {
	totals := newRoomData()
	for _, room := range p.closedRooms {
		totals.append(room)
	}
	totals.Name = totalsEntryName
	return totals
}

func (p *FlatParser) String() string {
	// if this sorting is done at room close (the closeRoom method) instead of here,
	// then it increases overall cpu usage with 90%
	sort.Slice(p.closedRooms, func(i, j int) bool {
		return p.closedRooms[i].Name < p.closedRooms[j].Name
	})
	roomStrings := []string{p.totals("total").String()}
	for _, room := range p.closedRooms {
		roomStrings = append(roomStrings, room.String())
	}
	return strings.Join(roomStrings, "\n")
}

func (p *FlatParser) HasOpenRooms() bool {
	return len(p.OpenRooms) > 0
}
