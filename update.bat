@echo off
d:
rmdir /S /Q d:\md\pti\pr
cd \md\pti\
git clone ivangit@192.168.1.12:/home/ivangit/pr
cd \combo\data
copy \md\pti\pr\DU\DU.xml >NUL
copy \md\pti\pr\Baz1\Baz1.xml >NUL
copy \md\pti\pr\Baz1\SBz1DU.xml >NUL
copy \md\pti\pr\Baz2\Baz2.xml >NUL
copy \md\pti\pr\Baz2\SBz2DU.xml >NUL
copy \md\pti\pr\scm\scm.xml >NUL
copy \md\pti\pr\scm\SBz1.xml >NUL
copy \md\pti\pr\scm\SBz2.xml >NUL
copy \md\pti\pr\scm\SDu.xml >NUL
echo All updated!