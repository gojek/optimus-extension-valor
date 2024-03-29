name="valor"
dir_output="dist"
os_list=("darwin" "linux" "windows")
arch_list=("amd64" "arm64")
if [ -z ${tag} ]
then
    tag="latest"
fi

for os in ${os_list[*]}; do
    for arch in ${arch_list[*]}; do
        file_name="${name}_${tag}_${os}-${arch}"
        path="${dir_output}/${file_name}"
        GOOS=$os GOARCH=$arch go build -trimpath -o ${path}
    done
done
