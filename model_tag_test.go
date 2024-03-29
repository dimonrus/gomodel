package gomodel

import (
	"testing"
)

func TestParseModelFiledTag(t *testing.T) {
	t.Run("all_in", func(t *testing.T) {
		tag := "col~created_at;seq;prk;frk~master.table(id,name);req;unq;cat;ign;"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
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
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_length", func(t *testing.T) {
		tag := "ca"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_tag", func(t *testing.T) {
		tag := "cac"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_frk", func(t *testing.T) {
		tag := "frk;aaa"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("wrong_col", func(t *testing.T) {
		tag := "col;aaa"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if field.Column != "" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("good_col", func(t *testing.T) {
		tag := "col~some_name;dat;uat"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if field.Column != "some_name" {
			t.Fatal("Wrong parser column name")
		}
	})
	t.Run("updated_at", func(t *testing.T) {
		tag := "col~updated_at;uat;"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if field.IsDeletedAt {
			t.Fatal("Wrong parser column id deleted at")
		}
	})
	t.Run("uat_dat_arr", func(t *testing.T) {
		tag := "col~data;uat;arr;dat;"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		if !field.IsDeletedAt {
			t.Fatal("Wrong parser column deleted at")
		}
		if !field.IsUpdatedAt {
			t.Fatal("Wrong parser column updated at")
		}
		if !field.IsArray {
			t.Fatal("Wrong parser column is array")
		}
		if len(tag) != len(field.String()) {
			t.Fatal("wrong compile")
		}
	})

}

func BenchmarkParseModelFiledTag(b *testing.B) {
	b.Run("all", func(b *testing.B) {
		tag := "col~created_at;seq;sys;prk;frk~master.table(id,name);req;unq;cat;"
		var field ModelFiledTag
		for i := 0; i < b.N; i++ {
			ParseModelFiledTag(tag, &field)
		}
		b.ReportAllocs()
	})

	b.Run("string", func(b *testing.B) {
		tag := "col~created_at;seq;sys;prk;frk~master.table(id,name);req;unq;cat;"
		var field ModelFiledTag
		ParseModelFiledTag(tag, &field)
		for i := 0; i < b.N; i++ {
			_ = field.String()
		}
		b.ReportAllocs()
	})

}
