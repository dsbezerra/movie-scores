@echo off

echo Build started.

REM Build for windows (caller arch)
echo Building for Windows...
go build -o movie-score-v2.exe
echo Building for Windows completed!

REM Build for Linux OS.
echo Building for Linux...
env GOOS=linux go build -o movie-score-v2-linux
echo Building for Linux completed!

echo Moving binaries to amenic-api/worker/bin...
REM Delete previous builds in bin folder of amenic-api/worker aka amenic-go.
rm -f ../amenic-api/worker/bin/movie-score-v2 && rm -f ../amenic-api/worker/bin/movie-score-v2.exe

REM Movie all new builds to the bin folder of amenic-api/worker aka amenic-go.
mv ./movie-score-v2.exe ../amenic-api/worker/bin && mv ./movie-score-v2-linux ../amenic-api/worker/bin

REM Script completed.
echo Move completed.
echo Build finished.
