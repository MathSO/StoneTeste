# Teste Stone
Repositório com minha solução ao teste técnico proposto pela empresa Stone

# Abstração
No enunciado pedia que baixassemos alguns dias de dados da plataforma B3 e fizessemos algumas agregações com elas podendo utlizar qualquer coisa a escolha do candidato, cada arquivo como observado, contem aproximadamente em média 10 milhoes de registro tendo assim que ser o mais otimizado possível. As agregações que eram solicitadas era o maior valor de venda/compra e o valor do dia com maior volume de negociações de um token podendo também ter um filtro por data até a data atual.

## Abordagem
Acredito que o primeiro desafio que encontrei foi como faria para subir todos os arquivos num banco de dados de uma maneira otimizada, tentei por um workbench (DBeaver) mas devido ao alto tempo que demorava e baixo desempenho resolvi criar o meu próprio script, assim poderia controlar melhor alguns aspectos como por exemplo paralelismo, no projeto que se encontra na pasta `go/importer` que contem um arquivo unico `import-csv.go` conta um pouco melhor como fiz e resolvi esse problema.

Após os arquivos devidamente importados e os dados no banco resolvi brincar um pouco com os dados para me familiarizar e entender realmente o que era proposto, assim fiz algumas queries, analizei seus desempenhos e resolvi criar alguns index para tornar-las mais rapidas, chegando assim a um resultado óitimo que também está dentro do arquivo `importer`, ou seja já é criado o index após a importação.

Resolvido o problema da query foi só fazer um servidor http para que pudesse chama-lá e distribuir da forma que foi apresentado pelo enunciado do problema. Esse projeto está dentro da pasta `go/server`, onde basicamente utilizei alguns frameworks e efetivamente resolvi o desafio.

Com o resultado satisfatorio em mãos criei uma imagem docker e docker compose para auxilixar a subir tudo, nele estão apenas o banco e server, caso queira importar arquivos deve-se rodar o projeto do importer com o comando `go run` ou construindo o executável e executando, lembrando que os arquivos com os dados devem estar dentro da pasta csv que se encontra dentro do projeto.

### Pastas e arquivos
* docker - Nessa pasta contem arquivos de imagens docker e tudo que precisa para a construção das imagens
* go - Aqui contém os códigos que montei para auxuliar e servir no problema, dividido em 2 subpastas, `importer` (projeto cru para importar qualquer arquivo csv que for colocado dentro da pasta `csv`) e `server` (projeto mais polido onde utilizo alguns frameworks para poder servir a resposta por http)
* postgres-data - Essa pasta só ira existir após a primeira vez que o docker compose for executado, contém o volume com as informações persistentes do postgres
* docker-compose.yaml - Arquivo para auxiliar a servir o banco de dados e o servidor
* READEME.md - Arquivo descritivo sobre abordagem e o projeto
* testes.sql - Arquivo com brainstorm de como cheguei a abordagem utilizada

#### Server
No servidor foi criado o endpoint `GET /info/:ticker` onde o ticker é o nome do codigo instrumento e também pode ser adicionar um parametro query `filter_date=yyyy-MM-dd` que faz o filto do dia dado até a ultima entrada.