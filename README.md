# clima-cep
 
API em Go que recebe um CEP, descobre a cidade pelo ViaCEP e devolve a temperatura atual (WeatherAPI) em Celsius, Fahrenheit e Kelvin.
 
Projeto do desafio "Clima por CEP no Google Cloud Run" do Go Expert.
 
## Em produção
 
https://clima-cep-260439605731.southamerica-east1.run.app/{cep}
 
Exemplo:
 
```
curl https://clima-cep-260439605731.southamerica-east1.run.app/60175047
```
 
Resposta:
 
```json
{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5}
```
 
Atenção: a URL sem o CEP no final retorna 404, porque a única rota é `GET /{cep}`.
 
## Respostas
 
| Situação | Status | Corpo |
|---|---|---|
| Sucesso | 200 | JSON com as três temperaturas |
| CEP sem 8 dígitos ou com caracteres inválidos | 422 | `invalid zipcode` |
| CEP com formato certo, mas que não existe | 404 | `can not find zipcode` |
 
As conversões usam `F = C * 1.8 + 32` e `K = C + 273`.
 
## Rodando localmente com Docker
 
Você precisa de uma chave gratuita da WeatherAPI (https://www.weatherapi.com/).
 
1. Copie o arquivo de exemplo e coloque sua chave:
```
cp .env.example .env
```
 
   Depois edite o `.env`:
 
```
WEATHER_API_KEY=sua_chave_aqui
```
 
2. Suba a aplicação:
```
docker compose up --build
```
 
3. Em outro terminal, faça uma consulta:
```
curl -i localhost:8080/60175047
```
 
## Rodando os testes
 
Com Docker:
 
```
docker compose run --rm test
```
 
Ou direto com Go (1.27 ou superior):
 
```
go test ./...
```
 
Os testes não acessam a internet nem precisam da chave, porque usam servidores e respostas simuladas.
 
## Estrutura
 
```
cmd/server/          ponto de entrada
internal/cep/        validação do CEP e cliente do ViaCEP
internal/weather/    cliente da WeatherAPI
internal/temperature/ conversões de temperatura
internal/handler/    handler HTTP
```
 
## Deploy
 
O serviço roda no Google Cloud Run (região southamerica-east1). Para publicar de novo, a partir da raiz do projeto:
 
```
gcloud run deploy clima-cep \
  --source . \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=sua_chave_aqui
```
 
