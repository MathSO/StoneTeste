select distinct datanegocio from negociacoes n;

select distinct codigoinstrumento from negociacoes n;

select count(1) from negociacoes n;

CREATE INDEX IF NOT EXISTS idx_codigoinstrumento ON negociacoes (codigoinstrumento);
CREATE INDEX IF NOT EXISTS idx_datanegocio ON negociacoes using BTREE (datanegocio desc nulls last);

CREATE INDEX IF NOT EXISTS idx_preconegocio_datanegocio ON negociacoes using BTREE (datanegocio desc nulls last, preconegocio asc nulls last);

CREATE INDEX IF NOT EXISTS idx_codigoinstrumento_preconegocio_datanegocio ON negociacoes using BTREE (codigoinstrumento, datanegocio desc nulls last, preconegocio asc nulls last);

DROP INDEX IF EXISTS idx_preconegocio_datanegocio;
DROP INDEX IF EXISTS idx_codigoinstrumento_preconegocio_datanegocio;

drop index idx_preconegocio_datanegocio;

explain analyze select sum(preconegocio) as sum_preco_negocio, count(1), max(preconegocio), datanegocio from negociacoes n  where codigoinstrumento = 'PETR4' and datanegocio >= DATE('2025-06-02') group by datanegocio order by sum_preco_negocio desc;

select sum(preconegocio) as sum_preco_negocio, count(1), max(preconegocio), datanegocio from negociacoes n  where codigoinstrumento = 'PETR4' and datanegocio >= DATE('2025-06-02') group by datanegocio order by sum_preco_negocio desc

explain analyze select count(1), max(preconegocio) from negociacoes n  where codigoinstrumento = 'PETR4' and datanegocio >= DATE('2025-06-02') order by sum_preco_negocio desc limit 1;

select max(preconegocio) from negociacoes n where datanegocio >= DATE('2025-06-02');
select count(1) from negociacoes n where datanegocio >= DATE('2025-06-02') and codigoinstrumento = 'AFHI11' group by n.datanegocio;
select sum(n.preconegocio), count(1) as sum_negocio from negociacoes n where datanegocio > DATE('2025-06-17') and codigoinstrumento = 'AFHI11'  order by sum_negocio desc limit 1;

select distinct datanegocio from negociacoes n;

select distinct codigoinstrumento from negociacoes n limit 1;


