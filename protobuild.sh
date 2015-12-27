
sourcePath="./src/GRts"

echo "-- Source Path --"
echo $sourcePath

echo "-- Create directories if necessary --"
if [[ ! -e $sourcePath/client/protocol ]]; then
    mkdir $sourcePath/client/protocol
fi
if [[ ! -e $sourcePath/server/protocol ]]; then
    mkdir $sourcePath/server/protocol
fi

echo "-- Remove previous protocol files --"
if ls $sourcePath/client/protocol/*.proto > /dev/null 2>&1; then
    rm $sourcePath/client/protocol/*.proto
fi
if ls $sourcePath/server/protocol/*.pb.go > /dev/null 2>&1; then
    rm $sourcePath/server/protocol/*.pb.go
fi

echo "-- Just copy them for the js client --"
cp $sourcePath/protocol/*.proto $sourcePath/client/protocol/

echo "-- Compile proto source to go for server --"
protoc --go_out=$sourcePath/server/protocol/  -I $sourcePath/protocol/ $sourcePath/protocol/*.proto

echo "-- Done --"
