files=`find $directory -type f -name "*.drawio"`
for i in $files
do
name=${i%.*}
drawio --export --format xml --output ./$name.png ./$name.drawio 
done


