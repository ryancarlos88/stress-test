To run this APP

clone the repository from 

git@github.com:ryancarlos88/stress-test.git

navigate to the cloned folder

run docker build -t stress-test . to build the docker image

run docker run stress-test -u https://httpstat.us/random/200-210 -r 45 -c9 as an example 

note.: this web site generates random status codes to see how they're presented at the final reports