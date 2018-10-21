FROM golang:1.10-windowsservercore-1803

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

RUN Invoke-WebRequest -UserAgent 'DockerCI' -outfile 7zsetup.exe "https://www.7-zip.org/a/7z1805-x64.exe" 
RUN $expectedSha1='A2183901F65B6789494AB712C54B6202A56CA5A0'; $fileSha1=(Get-FileHash -Path 7zsetup.exe -Algorithm SHA1).Hash; if ($fileSha1 -ne $expectedSha1) { Echo 'Bad SHA:' $fileSha1; Exit(1) }

RUN Invoke-WebRequest -UserAgent 'DockerCI' -outfile mingw64.7z http://downloads.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/7.3.0/threads-posix/seh/x86_64-7.3.0-release-posix-seh-rt_v5-rev0.7z
RUN $expectedSha1='0FCE15036400568BABD10D65B247E9576515DA2C'; $fileSha1=(Get-FileHash -Path mingw64.7z -Algorithm SHA1).Hash; if ($fileSha1 -ne $expectedSha1) { Echo 'Bad SHA:' $fileSha1; Exit(1) }

RUN start-process .\7zsetup -ArgumentList '/S /D=c:/7zip/' -Wait

RUN c:/7zip/7z x mingw64.7z -o'c:/'

RUN cmd.exe /C net users /ADD vcap /passwordreq:no /expires:never 
RUN cmd.exe /C runas /user:vcap whoami
RUN cmd.exe /C net accounts /maxpwage:UNLIMITED

RUN powershell.exe -Command Add-WindowsFeature Web-Webserver,Web-WebSockets,Web-WHC,Web-ASP,Web-ASP-Net45

RUN powershell.exe -command Set-ItemProperty -Path 'HKLM:\SYSTEM\CurrentControlSet\Services\dnscache' -Name Start -Value 4

RUN powershell.exe -command remove-windowsfeature -name 'windows-defender-features'

RUN powershell.exe -command ('MSDTC', 'TermService', 'dhcp', 'diagtrack', 'lmhosts', 'winrm', 'RemoteRegistry') | foreach { Set-Service -Name $_ -StartupType Disabled }
