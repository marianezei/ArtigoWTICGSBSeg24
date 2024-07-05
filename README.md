# Artigo WTICG SBSeg24
Título: Melhorias no Processo de Armazenamento de Dados em TPM para Gerenciamento de Integridade <br>
<br>
Resumo: Alguns dispositivos eletrônicos possuem soluções nativas para garantir sua integridade, um exemplo disso é o TPM (Trusted Platform Module), um chip dedicado a segurança. Nas máquinas virtuais, um vTPM(Virtual Trusted Platform Module) pode ser encontrado, este, quando ancorado com o TPM, pode usufruir da robustez de segurança que o TPM possui. Entretanto existe um obstáculo nessa estratégia, e é onde surge o objetivo deste trabalho. O vTPM, ao gerar múltiplas requisições ao TPM, pode gerar uma sobrecarga no chip e, para solucionar isso, o trabalho propõe a implementação de um escalonador de requisições, onde irá agrupar múltiplas chamadas e enviar somente uma ao TPM.
<br>

Os arquivos deste repositório foram desenvolvidos para serem executados apenas em máquinas que possuam um Trusted Platform Module (TPM). Além disso, é necessário que o usuário tenha permissões adequadas para realizar chamadas ao TPM. Certifique-se de que sua máquina possui um TPM e que você tem as permissões necessárias configuradas para garantir o funcionamento correto dos scripts e programas deste repositório.


Para executar os códigos deste repositório, certifique-se de ter o Go 1.20 instalado em seu ambiente de desenvolvimento. Além disso, este projeto depende do pacote `github.com/google/go-tpm` na versão `v0.9.1`. Você pode instalar essa dependência executando `go get github.com/google/go-tpm@v0.9.1`.
<br>


Antes de rodar o arquivo em Bash deste repositório, siga estas instruções:
- **Verifique as permissões de execução:** Certifique-se de que o arquivo tem permissão de execução. Você pode fazer isso com o comando <br> `chmod +x tpmTools_Concurrent.sh`.


Abaixo segue as descrições de cada arquivo:

- escalonador-concurrent.go : este script é a implementação do escalonador sugerido no artigo. Para executá-lo digite o comando <br>  ```go run escalonador-concurrent.go```.<br>
   Obs: o tempo inserido neste scritp não é necessário, o escalonador consegue lidar com qualquer volume de requisições.
  
- gotpm.go : este script apenas realiza as operações no TPM, foi desenvolvido para provar que a implementação do Escalonador traz melhorias. <br> Para executá-lo, digite o comando ```go run gotpm.go```.
  
- tpmTools_Concurrent.sh : este script foi implementado para a comparação do _tpm-tools_ com o escalonador utilizando _go-tpm_. <br> Para executá-lo ```./tpmTools_Concurrent.sh```.
<br>

Os demais arquivos de texto no repositório servem para a coleta de logs, leitura e escrita de funções, e um arquivo que é utilizado para a comparação do resultado obtido e o resultado esperado.
