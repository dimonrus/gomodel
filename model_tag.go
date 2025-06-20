package gomodel

import (
	"strings"
)

// ModelFiledTagList list of ModelFiledTag
type ModelFiledTagList []ModelFiledTag

// Len count of fields
func (l ModelFiledTagList) Len() int {
	return len(l)
}

// HasPrimary check if primary key exists
func (l ModelFiledTagList) HasPrimary() bool {
	for i := range l {
		if l[i].IsPrimaryKey {
			return true
		}
	}
	return false
}

// HasUnique check if unique key exists
func (l ModelFiledTagList) HasUnique() bool {
	for i := range l {
		if l[i].IsUnique {
			return true
		}
	}
	return false
}

// IsSoft check if model possible to soft delete
func (l ModelFiledTagList) IsSoft() bool {
	for i := range l {
		if l[i].IsDeletedAt {
			return true
		}
	}
	return false
}

// ModelFiledTag All possible model field tag properties
// tag must have 3 symbol lengths
type ModelFiledTag struct {
	// DB column name
	Column string `tag:"col"`
	// Foreign key definition
	ForeignKey string `tag:"frk"`
	// Has sequence
	IsSequence bool `tag:"seq"`
	// Is primary key
	IsPrimaryKey bool `tag:"prk"`
	// Is not null
	IsRequired bool `tag:"req"`
	// Is unique
	IsUnique bool `tag:"unq"`
	// Is created at column
	IsCreatedAt bool `tag:"cat"`
	// Is updated at column
	IsUpdatedAt bool `tag:"uat"`
	// Is deleted at column
	IsDeletedAt bool `tag:"dat"`
	// Is ignored column
	IsIgnored bool `tag:"ign"`
	// Is array value
	IsArray bool `tag:"arr"`
	// Interface to value
	Value any
	// If is zero
	IsZero bool
	// If is nil
	IsNil bool
	// Field position in struct
	Index int
}

// Clear tags
func (t *ModelFiledTag) Clear() {
	t.Column = ""
	t.ForeignKey = ""
	t.IsSequence = false
	t.IsPrimaryKey = false
	t.IsRequired = false
	t.IsUnique = false
	t.IsCreatedAt = false
	t.IsUpdatedAt = false
	t.IsDeletedAt = false
	t.IsIgnored = false
	t.IsArray = false
	t.Value = nil
	t.IsNil = false
	t.IsZero = false
	t.Index = 0
}

// Prepare string tag
func (t *ModelFiledTag) String() string {
	b := strings.Builder{}
	if t.Column != "" {
		b.WriteString("col~" + t.Column + ";")
	}
	if t.ForeignKey != "" {
		b.WriteString("frk~" + t.ForeignKey + ";")
	}
	if t.IsSequence {
		b.WriteString("seq;")
	}
	if t.IsPrimaryKey {
		b.WriteString("prk;")
	}
	if t.IsRequired {
		b.WriteString("req;")
	}
	if t.IsUnique {
		b.WriteString("unq;")
	}
	if t.IsCreatedAt {
		b.WriteString("cat;")
	}
	if t.IsUpdatedAt {
		b.WriteString("uat;")
	}
	if t.IsDeletedAt {
		b.WriteString("dat;")
	}
	if t.IsIgnored {
		b.WriteString("ign;")
	}
	if t.IsArray {
		b.WriteString("arr;")
	}
	return b.String()
}

// ParseModelFiledTag parse validation tag for rule and arguments
// Example
// db:"col~created_at;seq;sys;prk;frk~master.table(id,name);req;unq'"
func ParseModelFiledTag(tag string, field *ModelFiledTag) {
	if tag == "" || len(tag) < 3 {
		return
	}
	var indexStart, i int
	for i < len(tag) {
		if tag[i] == ';' || i == len(tag)-1 {
			switch tag[indexStart : indexStart+3] {
			case "seq":
				field.IsSequence = true
				i++
				indexStart = i
			case "prk":
				field.IsPrimaryKey = true
				i++
				indexStart = i
			case "req":
				field.IsRequired = true
				i++
				indexStart = i
			case "unq":
				field.IsUnique = true
				i++
				indexStart = i
			case "cat":
				field.IsCreatedAt = true
				i++
				indexStart = i
			case "uat":
				field.IsUpdatedAt = true
				i++
				indexStart = i
			case "dat":
				field.IsDeletedAt = true
				i++
				indexStart = i
			case "arr":
				field.IsArray = true
				i++
				indexStart = i
			case "ign":
				field.IsIgnored = true
				i++
				indexStart = i
			case "col":
				// Must be ~ according to format
				if tag[indexStart+3] != '~' {
					break
				}
				field.Column = tag[indexStart+4 : i]
				i++
				indexStart = i
			case "frk":
				// Must be ~ according to format
				if tag[indexStart+3] != '~' {
					break
				}
				field.ForeignKey = tag[indexStart+4 : i]
				i++
				indexStart = i
			}
		}
		i++
	}
	return
}
