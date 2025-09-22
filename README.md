# Extração de Dados do EPROC

## Descrição

Este serviço tem como objetivo extrair os dados do **AUTOR**, **RÉU** e **Valor da Causa** do sistema [eproc1g-consulta.tjsp.jus.br (EPROC)](https://eproc1g-consulta.tjsp.jus.br) e gerar um documento `.txt` com as seguintes informações:

- URL
- Número do Processo
- AUTOR
- RÉU
- Valor da Causa

## Funcionamento

O serviço realiza uma requisição HTTP para o sistema EPROC, captura a resposta, processa os dados relevantes do HTML e grava em um arquivo de texto estruturado.
