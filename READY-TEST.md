# MHRise NPLN Lab — Ready Test

Esta pasta contém uma preparação de laboratório baseada no projeto enviado.

## O que foi preparado

- certificado raiz exclusivo do laboratório: `certs/lab-root-ca.crt`
- chave da CA do laboratório: `certs/lab-root-ca.key`
- certificado do servidor NPLN assinado pela CA: `certs/server.crt`
- chave do servidor: `certs/server.key`
- override da CA usando o mecanismo já existente do Ryujinx `ssl`:
  `system/ssl/1054.der` (SecurityCommunicationRootCA)
- scripts de inicialização e reversão

## Teste rápido

Abra PowerShell na raiz e execute:

```powershell
.\scripts\run_ready_lab.ps1
```

O script:

1. usa o `server.exe` existente;
2. aponta `NEXTENDO_SERVER_IP` e `NEXTENDO_NAT_IP` para `127.0.0.1`;
3. cria um diretório de dados de laboratório com junctions para o save/configuração existentes;
4. instala temporariamente a CA do laboratório no certificado `1054`;
5. inicia o servidor;
6. inicia o Ryujinx-Nextendo com `--root-data-dir`.

## Depois do teste

Para remover o override da CA:

```powershell
.\scripts\restore_lab_ca.ps1
```

## O que observar

No Wireshark, capture `tcp.port == 443` na interface de loopback.

No log do Ryujinx procure o resultado do handshake NPLN.

No servidor procure chamadas gRPC.

### Resultado esperado do primeiro experimento

O objetivo imediato não é ainda ter todas as funções online do MHRise.

O primeiro marco é:

`ClientHello -> ServerHello/Certificate -> Finished -> HTTP/2/gRPC`

Se o erro `certificate verify failed` desaparecer, o próximo passo é capturar o primeiro RPC real e implementar a sequência exigida pelo cliente.

## Importante

A CA é somente para o laboratório local. Não é um certificado ou chave da Nintendo e não deve ser usada para contornar segurança em console real.
