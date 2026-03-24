//go:build wasm

package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// NotFoundPage adalah halaman 404 untuk route yang tidak dikenal.
type NotFoundPage struct {
	app.Compo
}

// Render menampilkan pesan 404 dengan tombol kembali ke beranda.
func (p *NotFoundPage) Render() app.UI {
	return app.Div().
		Style("min-height", "100vh").
		Style("background", "#0F0A1A").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("font-family", "'Inter', system-ui, sans-serif").
		Body(
			app.Div().
				Style("text-align", "center").
				Body(
					app.H1().
						Style("color", "#4D2975").
						Style("font-size", "72px").
						Style("font-weight", "800").
						Style("margin", "0 0 8px").
						Text("404"),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "16px").
						Style("margin", "0 0 24px").
						Text("Halaman tidak ditemukan."),
					app.A().
						Href("/").
						Style("background", "#4D2975").
						Style("color", "#E2D9F3").
						Style("text-decoration", "none").
						Style("padding", "10px 24px").
						Style("border-radius", "8px").
						Style("font-size", "14px").
						Style("font-weight", "600").
						Text("Kembali ke Beranda"),
				),
		)
}
