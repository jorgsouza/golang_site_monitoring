//hello.go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"bufio"
	"strconv"
	"github.com/fatih/color"
)

const monitoamentos = 3
const delay = 5

func exibeIntroducao()  {
	//declarar mome
	c := color.New(color.FgBlue, color.Bold)
	versao := 1.4
	c.Println("Monitoramento Sites Qlik" )
	fmt.Println("Este programa esta na versão", versao)
}

func main() {
	exibeIntroducao()
	for {
		exibeMenu()
		comando := leComando()

		//substiui o IF, com o switch não é necessário incluir o break, Apenas uma alternativa é executada por avaliação.
		switch comando {
		case 1:
			fmt.Println("")
			iniciarMonitoramento()
		case 2:
			exibeSubMenu()
		case 0:
			fmt.Println("Saindo do Programa.")
			os.Exit(0)
		default:
			fmt.Println("Comando Invalido")
			os.Exit(-1)
		}
	}
}

func exibeMenu(){
	//menu programa
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("0 - Sair do Programa")
}

func exibeSubMenu(){
	//Submenu programa
	g := color.New(color.FgHiGreen, color.Bold)
	g.Println("1 - Exibir todos os logs")
	g.Println("2 - Exibir apenas logs de erro")
	g.Println("3 - Voltar ao Menu Anterior \n")

	subcomando := leComandoSubmenu()
	switch subcomando {
	case 1:
		fmt.Println("")
		imprimeLogs()
	case 2:
		imprimeLogsErros()
	case 3:
		fmt.Println("Voltar ao Menu Anterior.")
		fmt.Println("")
	default:
		fmt.Println("Comando Invalido")
		os.Exit(-1)
	}

}

func leComando() int {
	//declarar opção programa
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("")

	return comandoLido
}

func leComandoSubmenu() int {
	//declarar opção programa
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("")

	return comandoLido
}

func iniciarMonitoramento()  {
	c := color.New(color.FgBlue, color.Bold)
	c.Println("Monitorando...")

	sites :=leSiteDoArquivo()

	for i :=0; i < monitoamentos ; i++ {
		for i, site := range sites {
			fmt.Println("Testando site", i, ":", site)
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
	fmt.Println("")
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registraLog(site, true)
	}else {
		r := color.New(color.FgHiRed, color.Bold)
		r.Println("Site:", site, "esta com problemas. Status Code:", resp.StatusCode)
		registraLog(site, false)
		registraLogErros(site, false)
	}
}

func leSiteDoArquivo()  []string{

	var sites []string

	arquivo, err := os.Open("sites.txt")
	if err != nil  {
		fmt.Println("Ocorreu um erro", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()

	return sites
}

func registraLog(site string, status bool)  {

	arquivo, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func registraLogErros(site string, status bool)  {

	arquivoErro, err := os.OpenFile("logErros.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	arquivoErro.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivoErro.Close()
}

func imprimeLogs()  {
	c := color.New(color.FgBlue, color.Bold)
	c.Println("Exibindo Logs...")
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(arquivo))
}

func imprimeLogsErros()  {
	r := color.New(color.FgHiRed, color.Bold)
	r.Println("Exibindo Logs de erro...")
	arquivoErro, err := ioutil.ReadFile("logErros.txt")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(arquivoErro))
}