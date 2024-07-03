# this script reads a file with a list of PCR values and extends them concurrently

#!/bin/bash
arquivo="input.txt"

if [ ! -f "$arquivo" ]; then
  echo "Erro: Arquivo '$arquivo' não encontrado."
  exit 1
fi

inicio=$(date +%s%N)
while IFS= read -r linha; do
  comando="tpm2_pcrextend 15:sha256=$linha"  # better choose an index that is not used by the system
  echo "Executando: $comando"
  $comando&
done < "$arquivo"

fim=$(date +%s%N)
tempo_total=$((fim - inicio))

echo "Tempo total de execução: $tempo_total nanosegundos"