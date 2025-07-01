import { useState } from "react"
import { useNavigate, Link } from "react-router-dom"
import Swal from "sweetalert2"


function Login() {
    const navigate = useNavigate()
    const [user_name, setUsername] = useState("")
    const [password, setPassword] = useState("")

    async function handleLogin(e){
        e.preventDefault()
        try {
            const requestBody = {user_name, password}
            console.log(requestBody);
            const respond= await fetch("http://localhost:8080/api/login",{
                method:"POST",
                headers:{"Content-Type":"application/json"},
                body:JSON.stringify(requestBody),
            }) 
            if(!respond.ok){
                const errorData = await respond.text();
                throw new Error(errorData || "login failed")
            } 

            const data= await respond.json();
            console.log("logged in : ",data);
            if(requestBody.user_name=='sarthak44'){
                navigate("/Admin", {
                    state: {
                        username: user_name,
                    }
                    });
            }
            else{
                navigate("/Dashboard", {
                    state: {
                        username: user_name,
                    }
                    });
            }
            
            
        } catch (error) {
            Swal.fire({
                icon: "error",
                title: "Oops...",
                text: error.message ||error.response.data.message
            });
        }
    }

    return (
        <>
            <div className=" flex  h-screen mx-auto items-center justify-center w-1/2" >
            <div className="w-1/2 bg-blue-50 p-8 rounded shadow-md items-center">
                    <form onSubmit={handleLogin} className="w-full">
                        <p>Welcome to CARPARK</p>
                        <div className="mb-4">
                            <label htmlFor="username" className="block mb-1 font-medium" >Username :</label>
                            <input  onChange={e => {setUsername(e.target.value)}} type="text" className="form-control bg-blue-200" id="username"/>
                        </div>
                        <div className="mb-3">
                            <label htmlFor="password" className="block mb-1 font-medium">Password :</label>
                            <input onChange={e => {setPassword(e.target.value)}} type="password" className="form-control bg-blue-200 border-black" id="password"/>
                        </div>
                        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition">LOG IN</button>
                        <p className="mt-4 text-sm">Don't have an account?<Link to={'/Signin'} className="text-blue-600 underline">Create an account</Link></p>

                    </form>
                </div>
            </div>
        </>
    )
}

export default Login