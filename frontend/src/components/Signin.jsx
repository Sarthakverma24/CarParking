import { useState } from "react"
import { useNavigate } from "react-router-dom"
import Swal from "sweetalert2"

function Signin() {
    const [Name, setName] = useState("")
    const [age, setage] = useState(18)
    const [username,setusername]=useState("")
    const [email, setEmail] = useState("")
    const [number,setnumber]=useState("")
    const [password, setPassword] = useState("")
    const navigate = useNavigate()

    async function handleSignin(e){
         e.preventDefault()
        const requestBody = {
            Name,
            age: Number(age), 
            username,
            email,
            number,
            password
        };  
        try{
            const res=await fetch("http://localhost:8080/api/signin",{
                method:"POST",
                headers:{"Content-Type":"application/json"},
                body:JSON.stringify(requestBody),
            })
            console.log(requestBody);
            console.log(typeof age);
            
            if(!res.ok){
                const error=await res.text();
                throw new Error(error|| "wrong details")
            }
            const data=await res.json();
            console.log(data)
            navigate('/')

        }
        catch(error){
            Swal.fire({
                icon: "error",
                title: "Oops...",
                text: error.message ||error.res.data.message
            }); 
        }
    }


    return (
        <>
            <div className=" flex  h-screen mx-auto items-center justify-center w-1/2" >
            <div className="w-1/2 bg-blue-50 p-8 rounded shadow-md items-center">
                    <form onSubmit={handleSignin} className="w-full">
                        <p>Welcome to CARPARK</p>
                        <div className="mb-4">
                            <label htmlFor="name" className="block mb-1 font-medium" >Name :</label>
                            <input  onChange={e => {setName(e.target.value)}} type="name" className="form-control bg-blue-200" id="name"/>
                        </div>

                        <div className="mb-4">
                            <label htmlFor="phonenumber" className="block mb-1 font-medium" >Phone Number :</label>
                            <input  onChange={e => {setnumber(e.target.value)}} type="phonenumber" className="form-control bg-blue-200" id="phonenumber"/>
                        </div>

                        <div className="mb-4">
                            <label htmlFor="email" className="block mb-1 font-medium" >Enter Email :</label>
                            <input  onChange={e => {setEmail(e.target.value)}} type="email" className="form-control bg-blue-200" id="email"/>
                        </div>

                        <div className="mb-4">
                            <label htmlFor="age" className="block mb-1 font-medium" >Age :</label>
                            <input  onChange={e => {setage(Number(e.target.value))}} type="number" className="form-control bg-blue-200" id="age"/>
                        </div>
        
                        <div className="mb-4">
                            <label htmlFor="username" className="block mb-1 font-medium" >UserName :</label>
                            <input  onChange={e => {setusername(e.target.value)}} type="username" className="form-control bg-blue-200" id="username"/>
                        </div>

                        <div className="mb-4">
                            <label htmlFor="password" className="block mb-1 font-medium">Password :</label>
                            <input onChange={e => {setPassword(e.target.value)}} type="password" className="form-control bg-blue-200 border-black" id="password"/>
                        </div>
                        
                        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition">Sign  IN</button>
                        
                    </form>
                </div>
            </div>
        </>
    )
}

export default Signin