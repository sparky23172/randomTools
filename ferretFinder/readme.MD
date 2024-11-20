# ferretFinder
Tools that I use to find specific things

# Tips

Jadx:
Entire file path is needed!!!

Works!
- jadx -d /home/kali/HTB/dontOverReact/jadxOutput /home/kali/HTB/dontOverReact/app-release.apk
- jadx -d `pwd`/out `pwd`/APKrypt.apk

Doesn't work
- jadx -d jadxOutput app-release.apk
