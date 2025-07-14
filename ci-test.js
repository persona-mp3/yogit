import { exec, execFile } from "node:child_process"
import fs from "node:fs/promises"

// so we want to able to write the content of the request into a file -> content.txt
// we can just write the html_url, author, and some/else
// and then the s1.sh can read the content.txt
// use the $("git clone --single-branch -b playground https://github.com/persona-mp3/yogit.git")
// 
const reqContent = {
    "time": "Jan 3, 1996, 2005",
    "user": "Yogit Charles",
    "device": "ASUS Racer of Gamers",
    "env": "NodeJS Eventloop timers",
}


async function ChildP() {
    try{
        const content = JSON.stringify(reqContent)
        /* 
        * fs.writeFile(filePath, content, options? )
        * options = {flag, objectEncoding, ...}
        */
        await fs.truncate("./content.txt", 0)
        await fs.writeFile("./content.txt", content, { flag: "a+"} )
        console.log("no errors occured, successfully wrote to content.txt...")
        
        console.log("changing permissons on s1.sh file")

        const fHandle = await fs.open("./w2.sh")
        await fHandle.chmod(0o777)

        console.log("\n---closing content file---\n")
        await fHandle.close()
    
    } catch(err) {
        console.log("an error occured trying to open content.txt", err)
        return
    }
    
    console.log("executing setup1.sh script")


    const child2 = exec(`bash s1.sh`, (err, stdout, stderr) => {
        if (err) {
            console.log("error occured in executing bash script\n",err)
            return
        }
        
        console.log("---script should be executing by now---\n")
        console.log(stdout)
    })
    
}

await ChildP()
