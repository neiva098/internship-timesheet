package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

func typeCredentials(ctx context.Context, user string, password string) bool {
	err := chromedp.Run(ctx,
		chromedp.SendKeys("#txtUsuario", user, chromedp.BySearch),
		chromedp.SendKeys("#txtSenha", password, chromedp.BySearch),
	)

	if err != nil {
		log.Fatal(err)
	}

	return true
}

func logIn(ctx context.Context, user string, password string) bool {
	chromedp.Run(ctx,
		chromedp.Navigate(`http://vtpro.seidorbrasil.com.br/`),
	)

	typeCredentials(ctx, user, password)

	chromedp.Run(ctx,
		chromedp.Click(`#bntAcessar`),
	)

	return true
}

func goToLancamentos(ctx context.Context) bool {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#ctl00_MenuPrincipaln1 > table > tbody > tr > td > a"),

		chromedp.Navigate("http://vtpro.seidorbrasil.com.br/Conteudo/TimeSheet/Lancamento.aspx"),
	)

	if err != nil {
		log.Fatal(err)
	}

	return true
}

func getAnoVtPro(ctx context.Context) string {
	var anoDoVtPro string

	const anoSelector = "#ctl00_cphPrincipal_TimeSheeet_lblMes"

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(anoSelector),

		chromedp.TextContent(anoSelector, &anoDoVtPro),
	)

	if err != nil {
		log.Fatal(err)
	}

	return anoDoVtPro
}

func initializeChrome() (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	ctx, cancel = context.WithTimeout(ctx, 9*time.Second)

	return ctx, cancel
}

func cakeRecipe(ctx context.Context, user string, password string) string {
	if user == "" || password == "" {
		log.Fatal("Não encontrei usuário ou senha nas variáveis de ambiente")
	}

	logIn(ctx, user, password)

	goToLancamentos(ctx)

	return getAnoVtPro(ctx)
}

func main() {
	godotenv.Load(".env")

	chrome, _ := initializeChrome()

	anoDoVtPro := cakeRecipe(chrome, os.Getenv("VTPRO_USER"), os.Getenv("VTPRO_PASSWORD"))

	fmt.Println("O mes no vtPro é ", anoDoVtPro)
}
