# go-rate-limiter

### Objetivo
Desenvolver um limitador de taxa em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

### Descrição
O objetivo deste desafio é criar um limitador de taxa em Go que possa ser usado para controlar o tráfego de requisições para um serviço web. O limitador de taxa deve ser capaz de limitar o número de requisições com base em dois critérios:

- **Endereço IP:** O limitador de taxa deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
- **Token de Acesso:** O limitador de taxa também deve ser capaz de limitar requisições com base em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser fornecido no cabeçalho no seguinte formato:
  `API_KEY: <TOKEN>`
  As configurações de limite do token devem substituir as do IP. Por exemplo, se o limite por IP for de 10 req/s e o de um token específico for de 100 req/s, o limitador de taxa deve usar as informações do token.

### Requisitos
- O limitador de taxa deve ser capaz de funcionar como um middleware que é injetado no servidor web.
- O limitador de taxa deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O limitador de taxa deve ter a opção de escolher o tempo de bloqueio para o IP ou Token se o número de requisições tiver sido excedido.
- As configurações de limite devem ser feitas por meio de variáveis de ambiente ou em um arquivo ".env" na pasta raiz.
- Deve ser possível configurar o limitador de taxa tanto para limitação de IP quanto de token de acesso.
- O sistema deve responder adequadamente quando o limite for excedido:
    - Código HTTP: 429
    - Mensagem: você atingiu o número máximo de requisições ou ações permitidas dentro de um determinado intervalo de tempo
- Todas as informações do limitador devem ser armazenadas e consultadas a partir de um banco de dados Redis. Você pode usar o docker-compose para iniciar o Redis.
- Crie uma "estratégia" que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limitador deve ser separada do middleware.

### Exemplos
- **Limitação de IP:** Suponha que o limitador de taxa esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
- **Limitação de Token:** Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições dentro desse intervalo, a décima primeira deve ser bloqueada.
  Em ambos os casos acima, as requisições subsequentes só podem ser feitas após o tempo total de expiração ter passado. Por exemplo, se o tempo de expiração for de 5 minutos, um IP específico só poderá fazer novas requisições após os 5 minutos terem passado.

### Dicas
- Teste seu limitador de taxa em diferentes condições de carga para garantir que funcione conforme o esperado em situações de alto tráfego.

### Entrega
- O código-fonte completo da implementação.
- Documentação explicando como o limitador de taxa funciona e como pode ser configurado.
- Testes automatizados demonstrando a eficácia e robustez do limitador de taxa.
- Use docker/docker-compose para que possamos testar sua aplicação.
- O servidor web deve responder na porta 8080.

---

## Documentação do Projeto

### Introdução
Este projeto utiliza um servidor Gin para fornecer serviços HTTP. Ele também faz uso de um sistema de controle de taxa de solicitações baseado em Redis ou em memória para evitar abusos.

### Configuração

Para executar este projeto corretamente, é necessário preencher as seguintes variáveis de ambiente:

- **PORT**: A porta em que o servidor estará escutando. Exemplo: `8080`.
- **RATE_LIMITER_STRATEGY**: A estratégia de controle de taxa de solicitações a ser usada. Exemplo: `redis`. Por padrão será utilizado um controle em memória.
- **RATE_LIMITER_IP_MAX_REQUESTS**: O número máximo de solicitações permitidas por IP dentro do intervalo de tempo especificado. Exemplo: `5`.
- **RATE_LIMITER_TOKEN_MAX_REQUESTS**: O número máximo de solicitações permitidas por token (`API_KEY`) dentro do intervalo de tempo especificado. Exemplo: `10`.
- **RATE_LIMITER_TIME_WINDOW_MILLISECONDS**: A janela de tempo em milissegundos na qual as solicitações são contadas para aplicação do limite. Exemplo: `10000` para 10 segundos.
- **RATE_LIMITER_BLOCKING_TIME_MILLISECONDS**: O tempo em milissegundos que um usuário será bloqueado após exceder o limite de solicitações. Exemplo: `20000` para 20 segundos.
- **REDIS_ADDR**: O endereço e porta do servidor Redis. Exemplo: `localhost:6379`.
- **REDIS_PASSWORD**: A senha, se necessário, para acessar o servidor Redis.

Certifique-se de que todas essas variáveis de ambiente estão devidamente configuradas antes de iniciar o projeto.

### Notas Adicionais
- Certifique-se de que o servidor Redis está em execução e acessível antes de iniciar o projeto.
- O servidor Gin e o servidor Redis devem estar disponíveis na rede para que o projeto funcione corretamente.

### Como executar?

#### Ambiente Dev
Altere o arquivo .env com os seguintes valores:

```
DOCKERFILE=Dockerfile.dev
IS_DEV=true
```

Execute o seguinte comando através do Docker Compose:

```shell
$ docker compose up -d
```

Conecte-se no container **ratelimit** e execute o serviço:

```shell
$ docker compose exec ratelimit sh
$ go run cmd/server/main.go
```

Dica: Utilize a extensão **Remote Development** no **VSCode** para realizar um ´Attach to running container´.

O serviço iniciará na porta 8080.
Para facilitar, utilize os arquivos **http** disponíveis no diretório **api**.

### Ambiente Produção
Altere o arquivo .env com os seguintes valores:

```
DOCKERFILE=Dockerfile.prod
IS_DEV=false
```

Execute o seguinte comando através do Docker Compose:

```shell
$ docker compose up --build
```

Os contêiner **ratelimit** estará pronto para uso, e você poderá realizar as chamadas HTTP.
O serviço iniciará na porta 8080.
Para facilitar, utilize os arquivos **http** disponíveis no diretório **api**.

### Testes unitários

Foi utilizado o pacote [gotestsum](https://github.com/gotestyourself/gotestsum)

```shell
$ gotestsum --format=short -- -coverprofile=coverage.out ./...
$ go tool cover -html=coverage.out -o coverage.html
```

Após a geração, abra o arquivo coverage.html para verificar a cobertura que deverá ultrapassar 90%.

### Testes de carga com rate limiting

Foi utilizado o pacote [bombardier](https://github.com/codesenberg/bombardier)

Exemplo com 5 conexões concorrentes, com duração de 10 segundos. E API com as seguintes configurações:

- RATE_LIMITER_IP_MAX_REQUESTS=5
- RATE_LIMITER_TIME_WINDOW_MILISECONDS=10000 (10s)
- RATE_LIMITER_BLOCKING_TIME_MILLISECONDS=2000 (2s)

```shell
$ bombardier -c 5 -d 10s http://localhost:8080/
```

Resultados:
```
Bombarding http://localhost:8080/ for 10s using 5 connection(s)
[=============================================================] 10s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec      2125.92     160.86    2641.69
  Latency        2.34ms   221.89us     6.60ms
  HTTP codes:
    1xx - 0, 2xx - 25, 3xx - 0, 4xx - 21240, 5xx - 0
    others - 0
  Throughput:   678.53KB/s
```

Exemplo com 5 conexões concorrentes, com duração de 10 segundos. E API com as seguintes configurações:

- RATE_LIMITER_TOKEN_MAX_REQUESTS=10
- RATE_LIMITER_TIME_WINDOW_MILISECONDS=10000 (10s)
- RATE_LIMITER_BLOCKING_TIME_MILLISECONDS=2000 (2s)

```shell
$ bombardier -c 5 -d 10s -H "API_KEY:123456" http://localhost:8080/
```

Resultados:
```
Bombarding http://localhost:8080/ for 10s using 5 connection(s)
[=============================================================================] 10s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec      2124.00     148.31    2583.94
  Latency        2.34ms   208.73us     5.00ms
  HTTP codes:
    1xx - 0, 2xx - 50, 3xx - 0, 4xx - 21196, 5xx - 0
    others - 0
  Throughput:   712.87KB/s
```