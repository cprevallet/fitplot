if [[ -z "$1" ]]; then
        echo "Usage startgui <Fit or TCX filename>"
else
        export GDK_BACKEND=x11
        ../../dumpfile $1 > test.dat
        ../../gui
        rm test.dat
fi
