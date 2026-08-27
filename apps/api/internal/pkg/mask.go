package pkg

import (
	"fmt"
	"regexp"
	"strings"
)

// Oculta sobrenomes de pessoas físicas (ex: "LUCAS SOARES DE FRANCA" -> "LUCAS S. D. F.")
func MascararNomePessoaFisica(nome string) string {
	partes := strings.Fields(strings.TrimSpace(nome))
	if len(partes) <= 1 {
		return nome
	}

	primeiroNome := partes[0]
	var inicais []string

	for _, p := range partes[1:] {
		if len(p) > 0 {
			inicais = append(inicais, string([]rune(p)[0])+".")
		}
	}

	return fmt.Sprintf("%s %s", primeiroNome, strings.Join(inicais, " "))
}

// Oculta parte do telefone de pessoas físicas (ex: "81988887777" -> "(81) *****-7777")
func MascararTelefone(tel string) string {
	apenasNumeros := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, tel)

	if len(apenasNumeros) < 8 {
		return "*****"
	}

	if len(apenasNumeros) >= 10 {
		ddd := apenasNumeros[:2]
		ultimos := apenasNumeros[len(apenasNumeros)-4:]
		return fmt.Sprintf("(%s) *****-%s", ddd, ultimos)
	}

	ultimos := apenasNumeros[len(apenasNumeros)-4:]
	return fmt.Sprintf("*****-%s", ultimos)
}

// MascararDocumento oculta dados do CPF se o documento tiver até 11 dígitos numéricos
// Mantém CNPJs (14 dígitos) e documentos de PJ inalterados.
// Ex: "12345678901" ou "123.456.789-01" -> "123.***.***-01"
func MascararDocumento(doc string) string {
	re := regexp.MustCompile(`\D`)
	numeros := re.ReplaceAllString(doc, "")

	if strings.TrimSpace(numeros) == "" {
		return doc
	}

	// Se possui até 11 dígitos, é tratado como CPF (Pessoa Física)
	if len(numeros) <= 11 {
		// Preenche com zeros à esquerda caso venha sem formatação e menor que 11
		cpfFormatado := fmt.Sprintf("%011s", numeros)
		blocoInicio := cpfFormatado[:3]
		blocoFim := cpfFormatado[len(cpfFormatado)-2:]

		return fmt.Sprintf("%s.***.***-%s", blocoInicio, blocoFim)
	}

	// CNPJ (14 dígitos) permanece visível sem mascaramento
	return doc
}
