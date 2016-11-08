#rsync -r -z -e "ssh -p 3145" dist.tar.gz remco@m.oornick.com:/home/remco/
rsync -r -z -e "ssh -p 3145" dist.tar.gz remco@127.0.0.1:/home/remco/


#git subtree split --prefix dist -b dist
#git push dewt dist:master
