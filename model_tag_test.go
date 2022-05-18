package gomodel

import (
	"testing"
)

func TestParseModelFiledTag(t *testing.T) {
	t.Run("all_in", func(t *testing.T) {
		tag := "col~created_at;seq;prk;frk~master.table(id,name);req;unq;cat;ign;"
		field := ParseModelFiledTag(tag)
		if field.Column != "created_at" {
			t.Fatal("Wrong parser column name")
		}
		if field.ForeignKey != "master.table(id,name)" {
			t.Fatal("Wrong parser fk")
		}
		if !field.IsRequired {
			t.Fatal("Wrong IsRequired")
		}
		if !field.IsCreatedAt {
			t.Fatal("Wrong is created at")
		}
		if !field.IsUnique {
			t.Fatal("Wrong IsUnique")
		}
		if !field.IsPrimaryKey {
			t.Fatal("Wrong IsPrimaryKey")
		}
		if !field.IsSequence {
			t.Fatal("Wrong IsSequence")
		}
		if !field.IsIgnored {
			t.Fatal("Wrong IsIgnored")
		}
		if len(tag) != len(field.String()) {
			t.Log("wrong length in string method")
		}
	})
	t.Run("empty", func(t *testing.T) {
		tag := ""
		field := ParseModelFiledTag(tag)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_length", func(t *testing.T) {
		tag := "ca"
		field := ParseModelFiledTag(tag)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_tag", func(t *testing.T) {
		tag := "cac"
		field := ParseModelFiledTag(tag)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_frk", func(t *testing.T) {
		tag := "frk;aaa"
		field := ParseModelFiledTag(tag)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_col", func(t *testing.T) {
		tag := "col;aaa"
		field := ParseModelFiledTag(tag)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("good_col", func(t *testing.T) {
		tag := "col~some_name;dat;uat"
		field := ParseModelFiledTag(tag)
		if field.Column != "some_name" {
			t.Fatal("Wrong parser column name")
		}
	})

}

func BenchmarkParseModelFiledTag(b *testing.B) {
	b.Run("all", func(b *testing.B) {
		tag := "col~created_at;seq;sys;prk;frk~master.table(id,name);req;unq;cat;"
		for i := 0; i < b.N; i++ {
			_ = ParseModelFiledTag(tag)
		}
		b.ReportAllocs()
	})

	b.Run("string", func(b *testing.B) {
		tag := "col~created_at;seq;sys;prk;frk~master.table(id,name);req;unq;cat;"
		field := ParseModelFiledTag(tag)
		for i := 0; i < b.N; i++ {
			_ = field.String()
		}
		b.ReportAllocs()
	})

}
