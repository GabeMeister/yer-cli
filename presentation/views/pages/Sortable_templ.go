// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package presentation_views_pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	components "GabeMeister/yer-cli/presentation/views/components"
	"GabeMeister/yer-cli/presentation/views/components/AnalyzeManuallyPage"
)

func Sortable(allEngineers []string, duplicateEngineers []string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div data-templ-id=\"sortable\"><form id=\"shared-form\" class=\"flex gap-2\" hx-post=\"/items\" hx-trigger=\"submit\" hx-target=\"body\" hx-swap=\"innerHTML\"><div class=\"m-4\"><span>Filter: </span> <input type=\"text\" name=\"filter-text\" id=\"filter-text\" class=\"ml-2 p-1 text-black\" hx-post=\"/search-engineers\" hx-trigger=\"input changed delay:500ms, keyup=[key==&#39;Enter&#39;], load\" hx-target=\"#left\" hx-swap=\"outerHTML\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = AnalyzeManuallyPage.AllEngineersList(allEngineers).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"right\" class=\"sortable w-96 h-96 p-6 bg-gray-700\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, engineer := range duplicateEngineers {
			templ_7745c5c3_Err = components.DraggableText(engineer).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input type=\"hidden\" name=\"ignore\" class=\"h-0 w-0\"> <input class=\"ignore-input\" type=\"hidden\" name=\"right-form-items\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(AnalyzeManuallyPage.GetCombinedValue(duplicateEngineers))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `presentation/views/pages/Sortable.templ`, Line: 39, Col: 134}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
