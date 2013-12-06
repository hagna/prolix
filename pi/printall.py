import filepath, time, sys, subprocess, fcntl

now = time.time()
pdelta = now - int(sys.argv[1])



class Lock(object):
    def __init__(self, name):
        self.name = name


    def lock(self):
        self.fd = file(self.name, 'w')
        try:
            fcntl.flock(self.fd.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
        except IOError, e:
            if e.errno == errno.EAGAIN:
                return False
            raise
        return True


if not Lock("/tmp/printall.lock").lock():
    print "could not acquire lock"
    sys.exit(1)

printdir = filepath.FilePath("/home/pi/print")
printed = printdir.child("printed")
errordir = printdir.child("errors")

for f in [printdir, printed, errordir]:
    if not f.exists():
        f.makedirs()

def printit(fname):
    try:
        print fname
        subprocess.check_output(["lpr", fname.path]) 
        fname.moveTo(printed.child(fname.basename()))
    except Exception, e:
        print e
        fname.moveTo(errordir.child(fname.basename()))


fp = filepath.FilePath("print")
for i in fp.walk(descend=lambda a: False):
    if i.basename().startswith("."):
        continue
    if i.isdir():
        continue
    mtime = i.getModificationTime()
    if mtime < pdelta:
        # it changed more than a pdelta ago
        printit(i)

