#!/usr/bin/env bash
  
parse_yaml () {
  local prefix=$2
  local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
  sed -ne "s|^\($s\):|\1|" \
    -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
    -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
  awk -F$fs '{
    indent = length($1)/2;
    vname[indent] = $2;
    for (i in vname) {if (i > indent) {delete vname[i]}}
    if (length($3) > 0) {
      vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
      printf("%s%s%s=%s\n", "'$prefix'",vn, $2, $3);
    }
  }'
}
  
yml=`parse_yaml ".github/workflows/heavy-tests.yml"`

OUTER_IFS=$IFS
IFS=$'\n'
for LINE in $yml; do
  echo "line $LINE"
  OIFS=$IFS
  IFS='='
  parts=($LINE)
  echo "parts ${parts[@]}"
  printf -v "${parts[0]}" '%s' "${parts[1]}"
  IFS=$OIFS
done
IFS=$OUTER_IFS
  
# while IFS= read -r LINE; do
#   echo "line $LINE"
#   OIFS=$IFS
#   IFS='='
#   parts=($LINE)
#   echo ${parts[@]}
#   printf -v "${parts[0]}" '%s' "$(echo ${parts[1]} | tr '§' ' ')"
#   IFS=$OIFS
# done < .tmp
  
  
echo $jobs_golangci_name
echo $jobs_test_name
echo $jobs_contract_test_name
