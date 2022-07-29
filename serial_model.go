package gomodel

import "github.com/dimonrus/gosql"

// SerialSoftTable init table for soft control model
func SerialSoftTable(name string) *gosql.Table {
	table := gosql.CreateTable(name)
	gosql.TableModeler{SerialPrimaryKeyModifier, TimestampModifier, SoftModifier}.Prepare(table)
	return table
}

// BigSerialSoftTable init table for soft control model
func BigSerialSoftTable(name string) *gosql.Table {
	table := gosql.CreateTable(name)
	gosql.TableModeler{BigSerialPrimaryKeyModifier, TimestampModifier, SoftModifier}.Prepare(table)
	return table
}

// SerialTable init table for control model
func SerialTable(name string) *gosql.Table {
	table := gosql.CreateTable(name)
	gosql.TableModeler{SerialPrimaryKeyModifier, TimestampModifier}.Prepare(table)
	return table
}

// BigSerialTable init table for control model
func BigSerialTable(name string) *gosql.Table {
	table := gosql.CreateTable(name)
	gosql.TableModeler{BigSerialPrimaryKeyModifier, TimestampModifier}.Prepare(table)
	return table
}

// TimestampModifier create timestamps
func TimestampModifier(tb *gosql.Table) {
	tb.AddColumn("created_at").Type("TIMESTAMP WITH TIME ZONE").Constraint().NotNull().Default("localtimestamp")
	tb.AddColumn("updated_at").Type("TIMESTAMP WITH TIME ZONE")
}

// SoftModifier create soft deleted_at column
func SoftModifier(tb *gosql.Table) {
	tb.AddColumn("deleted_at").Type("TIMESTAMP WITH TIME ZONE")
}

// SerialPrimaryKeyModifier create id serial primary key column
func SerialPrimaryKeyModifier(tb *gosql.Table) {
	tb.AddColumn("id").Type("serial").Constraint().NotNull().PrimaryKey()
}

// BigSerialPrimaryKeyModifier create id bigserial primary key column
func BigSerialPrimaryKeyModifier(tb *gosql.Table) {
	tb.AddColumn("id").Type("bigserial").Constraint().NotNull().PrimaryKey()
}
