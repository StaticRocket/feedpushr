package test

import (
	"context"
	"testing"

	"github.com/ncarlier/feedpushr/autogen/app"
	"github.com/ncarlier/feedpushr/autogen/app/test"
	"github.com/ncarlier/feedpushr/pkg/assert"
	"github.com/ncarlier/feedpushr/pkg/controller"
	"github.com/ncarlier/feedpushr/pkg/filter"
)

func TestFilterCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	ctrl := controller.NewFilterController(srv, filter.NewChainFilter())
	ctx := context.Background()

	// CREATE
	tags := "test"
	payload := &app.CreateFilterPayload{
		Name: "title",
		Props: map[string]interface{}{
			"prefix": "[test]",
		},
		Tags: &tags,
	}
	_, f := test.CreateFilterCreated(t, ctx, srv, ctrl, payload)
	assert.Equal(t, "title", f.Name, "")
	assert.Equal(t, false, f.Enabled, "")
	assert.ContainsStr(t, "test", f.Tags, "")
	assert.Equal(t, "[test]", f.Props["prefix"], "")
	assert.Equal(t, uint64(0), f.Props["nbSuccess"], "")
	id := f.ID

	// GET
	_, f = test.GetFilterOK(t, ctx, srv, ctrl, id)
	assert.Equal(t, id, f.ID, "")
	assert.Equal(t, "title", f.Name, "")

	// FIND
	_, list := test.ListFilterOK(t, ctx, srv, ctrl)
	assert.True(t, len(list) > 0, "")
	item := list[len(list)-1]
	assert.Equal(t, id, item.ID, "")

	// UPDATE
	tags = "test,foo"
	update := &app.UpdateFilterPayload{
		Enabled: true,
		Tags:    &tags,
	}
	_, f = test.UpdateFilterOK(t, ctx, srv, ctrl, id, update)
	assert.Equal(t, id, f.ID, "")
	assert.Equal(t, "title", f.Name, "")
	assert.ContainsStr(t, "test", f.Tags, "")
	assert.ContainsStr(t, "foo", f.Tags, "")
	assert.Equal(t, true, f.Enabled, "")

	// DELETE
	test.DeleteFilterNoContent(t, ctx, srv, ctrl, id)

	// GET 404
	test.GetFilterNotFound(t, ctx, srv, ctrl, id)
}
