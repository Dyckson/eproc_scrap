package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	numerosProcessos := []string{
		"40203176820258260100",
		"40203176820258260100",
	}

	file, err := os.Create("resultado_eproc.txt")
	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	for _, numeroDoProcesso := range numerosProcessos {
		url := fmt.Sprintf("https://eproc1g-consulta.tjsp.jus.br/eproc/externo_controlador.php?acao=processo_seleciona_publica&acao_origem=tjsp@consulta_unificada_publica/consultar&acao_retorno=tjsp@consulta_unificada_publica/consultar&num_processo=%s&num_chave=&num_chave_documento=&hash=1f9a21bc27e2ceca99b68396fcd5af85", numeroDoProcesso)

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Erro ao criar requisição para %s: %v", numeroDoProcesso, err)
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Referer", "https://www.linkedin.com/")
		req.Header.Set("Cookie", "PHPSESSID=cs74vrn9akj3517ss8p9r95j32")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Erro ao buscar URL para %s: %v", numeroDoProcesso, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Erro HTTP %d ao acessar %s", resp.StatusCode, url)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("Erro ao parsear HTML para %s: %v", numeroDoProcesso, err)
			continue
		}

		// ====== CAPTURAR AUTOR E RÉU ======
		var autor, reu string
		doc.Find("#fldPartes table tr").Each(func(i int, s *goquery.Selection) {
			if i == 1 {
				cols := s.Find("td")
				if cols.Length() >= 2 {
					autor = strings.TrimSpace(cols.Eq(0).Text())
					reu = strings.TrimSpace(cols.Eq(1).Text())
				}
			}
		})

		// ====== CAPTURAR VALOR DA CAUSA ======
		var valor string
		doc.Find("#fldInformacoesAdicionais table tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
			tds := s.Find("td")
			for j := 0; j < tds.Length()-1; j++ {
				label := strings.TrimSpace(tds.Eq(j).Text())
				if strings.HasPrefix(label, "Valor da Causa") {
					valor = strings.TrimSpace(tds.Eq(j + 1).Text())
					return false
				}
			}
			return true
		})

		fmt.Fprintf(file, "URL: %s\n", url)
		fmt.Fprintf(file, "Número do Processo: %s\n", numeroDoProcesso)
		fmt.Fprintf(file, "AUTOR: %s\n", autor)
		fmt.Fprintf(file, "RÉU: %s\n", reu)
		fmt.Fprintf(file, "Valor da Causa: %s\n", valor)
		fmt.Fprintf(file, "-----------------------------\n")
	}

	fmt.Println("✅ Dados extraídos do EPROC e salvos em resultado_eproc.txt")
}
