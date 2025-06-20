// testes feitos para importação incluiam 21 arquivos em aproximadamente 15 minutos, criação dos index cada de 3 a 5 minuitos, a quantidade total de registros foram de 194.054.308 dos dias 20/05/2025 até 17/06/2025
// tempo total nesse fluxo nos testes foram em torno de 25 minutos
package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// Diretorio onde vai procurar por arquivos para importar, não necessároamente precisa ser um arquivo com extensão csv
const CSV_DIR_NAME = "csv"

// csvFiles retorna um array com o os nomes dos arquivos na pasta
func csvFiles(dirName string) ([]string, error) {
	d, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	var fName []string
	for _, f := range d {
		if !f.IsDir() {
			fName = append(fName, f.Name())
		}
	}

	return fName, nil
}

// openCSV abre e retorna o arquivo sob o nome especificado em um *csv.Reader, assim cada comando Read() retorna uma linha, o cabeçalho já foi lido
func openCSV(path string) (*csv.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	csvReader.FieldsPerRecord = len(header)
	return csvReader, nil
}

// readCSVRows le 'rCount' linhas do Reader especificado, retornando os valores já divididos, caso chege no final do arquivo as linhas restantes são retornadas, juntamente com o erro io.EOF
func readCSVRows(rCount int, file *csv.Reader) (rows [][]string, err error) {
	var line []string
	for range rCount {
		line, err = file.Read()

		if err != nil {
			return
		}

		rows = append(rows, line)
	}

	return
}

// connectPostgres retorna uma conexão com o banco de dados, a string de conexão está fixada, mas pode ser editada. Também é passado um int que é passado para a pool de conexões, fazendo com que possa ter várias conexões ativas ao mesmo tempo, possibilitando multithreading
func connectPostgres(maxConn int32) (db *pgxpool.Pool, err error) {
	cfg, err := pgxpool.ParseConfig("host=localhost port=5432 user=stone password=stone dbname=stone sslmode=disable")
	if err != nil {
		return
	}

	cfg.MaxConns = maxConn

	return pgxpool.NewWithConfig(context.Background(), cfg)
}

// createTable cria a tabela onde será importado os dados na conexão enviada
func createTable(conn *pgxpool.Pool) error {
	batch := &pgx.Batch{}
	batch.Queue(`CREATE TABLE IF NOT EXISTS negociacoes (
		data_referencia DATE,
		codigo_instrumento VARCHAR(20),
		acao_atualizacao INT,
		preco_negocio DECIMAL(18, 4),
		quantidade_negociada INT,
		hora_fechamento VARCHAR(20),
		codigo_identificador_negocio INT,
		tipo_sessao_pregao INT,
		data_negocio DATE,
		codigo_participante_comprador INT,
		codigo_participante_vendedor INT)`)

	err := conn.SendBatch(context.Background(), batch).Close()
	if err != nil {
		return err
	}

	return err
}

// insertRows insere as linhas que foram enviados utilizando protocolo de copia
func insertRows(rows [][]string, conn *pgxpool.Pool) error {
	_, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"negociacoes"},
		[]string{
			"data_referencia",
			"codigo_instrumento",
			"acao_atualizacao",
			"preco_negocio",
			"quantidade_negociada",
			"hora_fechamento",
			"codigo_identificador_negocio",
			"tipo_sessao_pregao",
			"data_negocio",
			"codigo_participante_comprador",
			"codigo_participante_vendedor",
		},
		pgx.CopyFromSlice(len(rows), func(i int) (row []any, err error) {
			for _, v := range rows[i] {
				v = strings.ReplaceAll(v, ",", ".")

				if v == "" {
					row = append(row, nil)
				} else {
					row = append(row, v)
				}
			}

			return
		}),
	)

	return err
}

// createIndex cira os indexes que serão utilizados na consulta pelo server
func createIndex(conn *pgxpool.Pool) error {
	defer fmt.Println("finished creating indexes")

	batch := &pgx.Batch{}

	// o index otimiza a consulta por codigo_instrumento, e ordena por data_negocio e preco_negocio, assim facilitando queries que buscem por um ticker e agrupem por data, bem como achar os minimos e máximos
	fmt.Println("creating index idx_codigoinstrumento_preconegocio_datanegocio")
	batch.Queue(`CREATE INDEX IF NOT EXISTS idx_codigoinstrumento_preconegocio_datanegocio ON negociacoes using BTREE (codigo_instrumento, data_negocio desc nulls last, preco_negocio asc nulls last)`)

	// o index otimiza consulta por codigo_instrumento
	fmt.Println("creating index idx_codigoinstrumento")
	batch.Queue(`CREATE INDEX IF NOT EXISTS idx_codigoinstrumento ON negociacoes (codigo_instrumento)`)

	err := conn.SendBatch(context.Background(), batch).Close()
	if err != nil {
		return err
	}

	return err
}

func listAllCodigoInstrumento(conn *pgxpool.Pool) error {
	file, err := os.Create("codigo_instrumento.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("listing all codigo instrumento")
	defer fmt.Println("listing completed")

	rows, err := conn.Query(context.Background(), `select distinct codigo_instrumento from negociacoes n`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var name string
		rows.Scan(&name)
		file.WriteString(name + "\n")
	}

	return nil
}

func main() {
	// Busca os asquivos
	files, err := csvFiles(CSV_DIR_NAME)
	if err != nil {
		panic(err)
	}
	totalFiles := len(files)

	// Conecta no banco
	conn, err := connectPostgres(int32(totalFiles))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Cria a tabela
	err = createTable(conn)
	if err != nil {
		panic(err)
	}

	// Cria a estrutura que para auxilixar nas goroutine
	wg := sync.WaitGroup{}
	wg.Add(totalFiles) // adiciona o total de arquivos ao grupo
	for _, fName := range files {
		// Abre uma goroutine por arquivo
		go func() {
			defer wg.Done() // finaliza um arquivo
			if fName == `.gitkeep` {
				return
			}

			fmt.Println("Begining file", fName)
			defer fmt.Println("Finished file", fName)

			// Abre o arquivo dessa rotina
			csvFile, err := openCSV(CSV_DIR_NAME + "/" + fName)
			if err != nil {
				panic(err)
			}

			block := 0
			for { // le o arquivo até o final e salva no banco
				block++
				rows, err := readCSVRows(100000, csvFile)

				if err != nil {
					if err == io.EOF && len(rows) > 0 {
						insertRows(rows, conn)
					}
					return
				}

				err = insertRows(rows, conn)
				if err != nil {
					fmt.Println("Aborting file", fName, "error on block", block)
					return
				}
			}
		}()
	}

	wg.Wait() // espera que todos os arquivos sejam finalizados

	err = createIndex(conn) // cria o index
	if err != nil {
		panic(err)
	}

	err = listAllCodigoInstrumento(conn) // lista os tickers lidos no arquivo codigo_instrumento.txt
	if err != nil {
		panic(err)
	}
}
