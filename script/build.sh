name="valor"
dir_output="dist"
os_list=("darwin" "linux")
arch_list=("amd64")
if [ -z ${tag} ]
then
    tag="latest"
fi

for os in ${os_list[*]}; do
    for arch in ${arch_list[*]}; do
        file_name="${name}_${tag}_${os}-${arch}"
        path="${dir_output}/${file_name}"
        GOOS=$os GOARCH=$arch go build -o ${path}
    done
done
