package src

import (
	"fmt"
	"strings"
)

type segment struct {
	start   int
	content string
}

func (s *segment) IsSame(seg *segment) bool {
	return s.start == seg.start && s.content == seg.content
}

func (s *segment) String() string {
	return fmt.Sprintf("[%d: '%s']", s.start, s.content)
}

func (s *segment) Overlaps(set LineSegments) LineSegments {
	return Overlaps(s.start, len(s.content), set)
}

func (s *segment) IsInSet(set LineSegments) bool {
	for _, seg := range set {
		if s.IsSame(seg) {
			return true
		}
	}
	return false
}

type LineSegments []*segment

// NewLineSegments is a constructor for LineSegments.
// It accepts initial segments as variadic params.
func NewLineSegments(segments ...*segment) LineSegments {
	ls := make(LineSegments, 0)
	if len(segments) > 0 {
		ls = append(ls, segments...)
	}
	return ls
}

func Split(line string) LineSegments {
	delimiters := map[string]struct{}{
		"|": {}, "\\": {}, "/": {}, "+": {}, "-": {},
	}

	segments := NewLineSegments()
	start := -1
	var foundFirstDelimiter bool

	for i, c := range line {
		if _, exists := delimiters[string(c)]; exists {
			if start >= 0 && foundFirstDelimiter && i > start {
				segment := &segment{start, strings.Trim(line[start:i], "+-")}
				segments = append(segments, segment)
			}
			foundFirstDelimiter = true
			start = i + 1
			continue
		}
		if start == -1 && foundFirstDelimiter {
			start = i
		}
	}

	return segments
}

// Overlaps return all the segments in set that overlap with the segment defined by start and length
func Overlaps(start, length int, set LineSegments) LineSegments {
	overlaps := NewLineSegments()

	if set == nil || len(set) == 0 {
		return overlaps
	}

	for _, segment := range set {
		if segment.start+len(segment.content) <= start || start+length <= segment.start {
			continue
		}
		overlaps = append(overlaps, segment)
	}

	return overlaps
}

// MultipleOverlaps finds the segments from set2 that overlap with any of the segments in set1.
// nonOverlappingSegments will be the rest.
func MultipleOverlaps(set1, set2 LineSegments) (overlappingSegments, nonOverlappingSegments LineSegments) {
	for _, s2 := range set2 {
		for _, s1 := range set1 {
			over := s1.Overlaps(NewLineSegments(s2))
			overlappingSegments = append(overlappingSegments, over...)
		}
	}
	nonOverlappingSegments = segmentsDiff(set2, overlappingSegments)
	return
}

// segmentsDiff will return all the segments in set1 that are not in set2:
// set1 minus set2
func segmentsDiff(set1, set2 LineSegments) (result LineSegments) {
	if len(set1) == 0 {
		return
	}
	if len(set2) == 0 {
		return set1
	}
	for _, s1 := range set1 {
		if !s1.IsInSet(set2) {
			result = append(result, s1)
		}
	}
	return
}

// sendData looks for title and chairs inside a string
func segmentData(str string) (*roomData, error) {
	var roomTitle string
	roomData := newRoomData()

	stillInsideTitle := false
	for _, c := range str {
		switch {
		case c == '(':
			stillInsideTitle = true
		case c == ')':
			stillInsideTitle = false
			roomData.Name = strings.TrimSpace(roomTitle)
		case stillInsideTitle:
			roomTitle += string(c)
		case c == 'W', c == 'P', c == 'S', c == 'C':
			roomData.Chairs[c]++
		case c == ' ':
		default:
			return nil, fmt.Errorf("strange character encountered: %c", c)
		}
	}

	if stillInsideTitle {
		return nil, fmt.Errorf("room title did not close. It starts with '%s'", roomTitle)
	}

	return roomData, nil
}
