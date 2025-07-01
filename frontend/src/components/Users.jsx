import {useState,useEffect} from 'react'
import { useNavigate} from "react-router-dom"
import Swal from "sweetalert2"; 
function Users() {
    const [users, setUsers] = useState([]);
    const navigate = useNavigate()
    useEffect(() => {
        fetchUsers();
      }, []);
    
      async function fetchUsers() {
        try {
          const response = await fetch("http://localhost:8080/api/Users");
          if (!response.ok) {
            throw new Error('Failed to fetch lots');
          }
          const data = await response.json();
          setUsers(data)
          
        } catch (error) {
          console.error("Error fetching lots:", error);
          Swal.fire({
            icon: "error",
            title: "Error",
            text: error.message
          });
        }
      }
    return (
        <>
            <div className="min-h-screen bg-blue-950">
                <div className="flex items-center justify-between bg-blue-900 p-4 w-full">
                    <h1 className="text-xl text-white font-semibold">Welcome Admin</h1>
                    <div className="flex gap-4">
                        <button className="px-4 py-2 bg-white rounded-lg shadow hover:-translate-y-1 hover:scale-110 hover:bg-white transition-colors"
                        onClick={()=>{navigate('/Admin')}}>
                            Home
                        </button>
                        <button className="px-4 py-2 bg-white rounded-lg shadow hover:-translate-y-1 hover:scale-110 hover:bg-white transition-colors"
                        onClick={()=>{navigate('/Users')}}>
                            Users
                        </button>
                        
                        <button className="px-4 py-2 bg-white rounded-lg shadow hover:-translate-y-1 hover:scale-110 hover:bg-white transition-colors"
                        onClick={()=>{navigate('/Summary')}}>
                            Summary
                        </button>
                        <button className="px-4 py-2 bg-white rounded-lg shadow hover:-translate-y-1 hover:scale-110 hover:bg-white transition-colors"
                        onClick={()=>{navigate('/')}}>
                            Logout
                        </button>
                    </div>
                    <div>
                    <h1 className="text-base underline cursor-pointer  hover:text-blue-700 transition-colors">
                        Edit Profile
                    </h1>
                    </div>
                </div>

                <div className="p-4 bg-blue-900 border border-black m-5 rounded-lg overflow-x-auto">
                    <table className="min-w-full text-white text-left border-collapse">
                        <thead>
                        <tr className="bg-blue-800 text-white">
                            <th className="px-4 py-2 border-b border-gray-600">Name</th>
                            <th className="px-4 py-2 border-b border-gray-600">Username</th>
                            <th className="px-4 py-2 border-b border-gray-600">Email</th>
                            <th className="px-4 py-2 border-b border-gray-600">Phone</th>
                            <th className="px-4 py-2 border-b border-gray-600">Age</th>
                        </tr>
                        </thead>
                        <tbody>
                        {users.length === 0 ? (
                            <tr>
                            <td colSpan="5" className="text-center py-4">No users found.</td>
                            </tr>
                        ) : (
                            users.map((user, index) => (
                            <tr key={index} className="hover:bg-blue-700 transition">
                                <td className="px-4 py-2 border-b border-gray-700">{user.name}</td>
                                <td className="px-4 py-2 border-b border-gray-700">{user.user_name}</td>
                                <td className="px-4 py-2 border-b border-gray-700">{user.gmail}</td>
                                <td className="px-4 py-2 border-b border-gray-700">{user.phone_no}</td>
                                <td className="px-4 py-2 border-b border-gray-700">{user.age}</td>
                            </tr>
                            ))
                        )}
                        </tbody>
                    </table>
                    </div>
            </div>
        </>
    )
}

export default Users
