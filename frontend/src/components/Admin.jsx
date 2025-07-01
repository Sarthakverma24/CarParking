import Swal from "sweetalert2";
import { useState, useEffect } from "react";
import { useLocation ,useNavigate} from "react-router-dom"

function Admin() {
  const navigate = useNavigate()
  const locations = useLocation();
  const { user_name } = locations.state || {};

  const [popup, setPopup] = useState(false);
  const[slotpop,setslotpop]=useState(false);
  const [lots, setLots] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedSlot, setSelectedSlot] = useState(null);

  
  const [location, setLocation] = useState("");
  const [address, setAddress] = useState("");
  const [pincode, setPincode] = useState("");
  const [price, setPrice] = useState("");
  const [spots, setSpots] = useState("");
  const [editingLot, setEditingLot] = useState(null);

  useEffect(() => {
    fetchLots();
    setSelectedSlot(slotpop);  
    setslotpop(true);       
  }, []);

  async function fetchLots() {
    setLoading(true);
    try {
      const response = await fetch("http://localhost:8080/api/lots");
      if (!response.ok) {
        throw new Error('Failed to fetch lots');
      }
      const data = await response.json();

      setLots(Array.isArray(data) ? data : []);
      console.log(data);
      
    } catch (error) {
      console.error("Error fetching lots:", error);
      Swal.fire({
        icon: "error",
        title: "Error",
        text: error.message
      });

      setLots([]);
    } finally {
      setLoading(false);
    }
  }


  async function addOrUpdateLot(e) {
    e.preventDefault();
    try {
      const method = editingLot ? "PUT" : "POST";
      const url = editingLot 
        ? `http://localhost:8080/api/lots/${editingLot.id}`
        : "http://localhost:8080/api/login/addlot";
      
      const requestBody = {
        ...(editingLot && { id: editingLot.id }),
        location,
        address,
        pincode,
        price: Number(price),
        spots: Number(spots)
      };
      
      const lotResponse = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestBody),
      });

      if (!lotResponse.ok) {
        const errorData = await lotResponse.text();
        throw new Error(errorData || "Lot operation failed");
      }

      const lotData = await lotResponse.json();
      const targetLotId = editingLot ? editingLot.id : lotData.lot_id;

      const spotsResponse = await fetch(
        `http://localhost:8080/api/lots/spots/${targetLotId}`
      );
      
      if (!spotsResponse.ok) {
        throw new Error("Failed to fetch spot data");
      }

      const spotsData = await spotsResponse.json();

      let successMessage = editingLot 
        ? "Lot updated successfully" 
        : "Lot added successfully";
      
      successMessage += `\nSpots: ${spotsData.length}`;
      
      console.log("Spots data:", spotsData);

      Swal.fire({
        icon: "success",
        title: "Success!",
        text: successMessage,
        footer: `${spotsData.length} spots retrieved`
      });
      
      setPopup(false);
      setEditingLot(null);
      fetchLots();
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Error",
        text: error.message
      });
    }
  }

  
  async function deleteSlot(e) {
      e.preventDefault();
      try {
        const parkingId = selectedSlot.parking_id;
        const url = `http://localhost:8080/api/slots/${parkingId}`;
        console.log("Deleting slot ID:", parkingId);

        const response = await fetch(url, {
          method: "DELETE",
        });

        if (!response.ok) {
          const errorData = await response.text();
          throw new Error(errorData || "Slot operation failed");
        }

        Swal.fire({
          icon: "success",
          title: "Success!",
          text: "Slot deleted successfully"
        });

        setslotpop(false);
        fetchLots();
      } catch (error) {
        Swal.fire({
          icon: "error",
          title: "Error",
          text: error.message
        });
      }
    }


  async function deleteLot(lotId) {
    try {
      const response = await fetch(`http://localhost:8080/api/lot/${lotId}`, {
        method: "DELETE"
      });
    
      
      if (!response.ok) {
        throw new Error('Failed to delete lot');
      }
      
      Swal.fire({
        icon: "success",
        title: "Deleted!",
        text: "Parking lot has been deleted."
      });
      
      fetchLots();
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Error",
        text: error.message
      });
    }
  }

  const openEditPopup = (lot) => {
    setEditingLot(lot);
    setLocation(lot.location || "");
    setAddress(lot.address || "");
    setPincode(lot.pincode || "");
    setPrice(lot.price || "");
    setSpots(lot.spots || "");
    setPopup(true);
  };

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

      <div className="p-4 bg-black flex justify-center border text-white  border-black m-5 rounded-lg">
        <h5 className="text-lg font-bold">Parking Lots</h5>
      </div>

      {/* Parking Lots Grid */}
      {loading ? (
        <div className="flex justify-center items-center h-64">
          <p className="text-white text-lg">Loading parking lots...</p>
        </div>
      ) : lots.length === 0 ? (
        <div className="flex justify-center items-center h-64">
          <p className="text-white text-lg">No parking lots found</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 p-4">
          {lots.map((lot) => (
            <div 
              key={lot._id}
              className="relative rounded-2xl border border-neutral-800 bg-neutral-900 overflow-hidden p-6 text-white shadow-2xl"
            >
              {/* Lot Header */}
              <div className="text-center mb-4">
                <h2 className="text-xl font-bold">
                  Parking ID: {lot.id || 'N/A' || 'N/A'}
                </h2>
                <p className="text-sm text-gray-400">{lot.location}</p>
              </div>
              
              {/* Edit/Delete Buttons */}
              <div className="flex justify-center gap-3 mb-4">
                <button 
                  onClick={() => openEditPopup(lot)}
                  className="px-3 py-1 bg-blue-500 rounded hover:bg-blue-600 transition"
                >
                  Edit
                </button>
                <button 
                  onClick={() =>{ deleteLot(lot.id);
                    console.log(lot._id);
                    
                  }}
                  className="px-3 py-1 bg-red-500 rounded hover:bg-red-600 transition"
                >
                  Delete
                </button>
              </div>
              
              {/* Parking Spots */}
              <div className="grid grid-cols-5 gap-2">
                {lot.spot_list && lot.spot_list.length > 0 ? (
                lot.spot_list.map((spot) => (
                  <button
                    key={spot.parking_id}
                    className={`w-8 h-8 rounded flex items-center justify-center text-xs transition ${
                      spot.status ? 'bg-green-600 hover:bg-green-500' : 'bg-red-600 hover:bg-red-500'
                    }`}
                    onClick={() => {
                      setSelectedSlot(spot);     
                      setslotpop(true);         
                    }}
                  >
                  </button>
                ))
              ) : (
                <p className="text-sm text-gray-400 col-span-5 text-center">No spots available</p>
              )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Add Lot Button */}
      <div className="fixed bottom-6 left-1/2 transform -translate-x-1/2">
        <button
          className="bg-red-500 hover:-translate-y-1 hover:scale-110 hover:bg-red-600 text-white font-semibold py-2 px-6 rounded-lg shadow-md transition"
          onClick={() => {
            setEditingLot(null);
            setLocation("");
            setAddress("");
            setPincode("");
            setPrice("");
            setSpots("");
            setPopup(true);
          }}
        >
          + Add Lot
        </button>
      </div>
      {/* View/delete slot Popup */}
      {slotpop && selectedSlot && (
        <div className="fixed inset-0 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl p-6 w-[400px] max-w-full shadow-xl relative border border-gray-300">
            
            {/* Close Button */}
            <button
              className="absolute top-3 right-4 text-gray-500 hover:text-gray-800 font-bold text-2xl"
              onClick={() => setslotpop(false)}
            >
              ×
            </button>

            <h2 className="text-xl font-bold text-gray-800 mb-6 text-center">
              View/Delete Parking Slot
            </h2>

            <form onSubmit={deleteSlot} className="space-y-4">

              {/* Slot ID (Read-only) */}
              <div className="flex items-center gap-2">
                <label htmlFor="slotId" className="text-black font-medium w-24">
                  Slot ID:
                </label>
                <input
                  type="text"
                  id="slotId"
                  value={selectedSlot.parking_id}
                  readOnly
                  className="form-control border border-black bg-gray-100 text-black px-2 py-1 flex-1 rounded"
                />
              </div>

              {/* Slot Status (Read-only) */}
              <div className="flex items-center gap-2">
                <label htmlFor="status" className="text-black font-medium w-24">
                  Status:
                </label>
                <input
                  type="text"
                  id="status"
                  value={selectedSlot.status ? "Occupied" : "Vacant"}
                  readOnly
                  className="form-control border border-black bg-gray-100 text-black px-2 py-1 flex-1 rounded"
                />
              </div>

              {/* Delete Button */}
              <div className="flex justify-center pt-4">
                <button
                  type="submit"
                  className="bg-red-600 text-white py-2 px-8 rounded hover:-translate-y-1 hover:scale-105 hover:bg-red-500 transition"
                >
                  Delete
                </button>
              </div>

            </form>
          </div>
        </div>
      )}
      {/* Add/Edit Lot Popup */}
      {popup && (
        <div className="fixed inset-0 flex items-center justify-center z-50">
            <div className="bg-white rounded-2xl p-6 w-[400px] max-w-full shadow-xl relative border border-gray-300">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-bold text-gray-800">
                {editingLot ? "Edit Parking Lot" : "New Parking Lot"}
              </h2>
              <button
                className="text-gray-500 hover:text-gray-800 font-bold text-xl"
                onClick={() => setPopup(false)}
              >
                ×
              </button>

            </div>
            <form onSubmit={addOrUpdateLot} className="w-full">
              <div className="mb-3 flex items-center gap-2">
                <label htmlFor="location" className="text-black font-medium w-24">Location:</label>
                <input
                  value={location}
                  onChange={e => setLocation(e.target.value)}
                  type="text"
                  className="form-control border border-black text-black px-2 py-1 flex-1 rounded"
                  id="location"
                  required
                />
              </div>
              <div className="mb-3 flex items-center gap-2">
                <label htmlFor="address" className="text-black font-medium w-24">Address:</label>
                <input
                  value={address}
                  onChange={e => setAddress(e.target.value)}
                  type="text"
                  className="flex-grow h-12 px-4 py-2 border border-black text-black text-lg rounded"
                  id="address"
                  required
                />
              </div>
              <div className="mb-3 flex items-center gap-2">
                <label htmlFor="pincode" className="text-black font-medium w-24">Pincode:</label>
                <input
                  value={pincode}
                  onChange={e => setPincode(e.target.value)}
                  type="number"
                  className="form-control border border-black text-black px-2 py-1 flex-1 rounded"
                  id="pincode"
                  required
                />
              </div>
              <div className="mb-3 flex items-center gap-2">
                <label htmlFor="price" className="text-black font-medium w-24">Price/hr:</label>
                <input
                  value={price}
                  onChange={e => setPrice(e.target.value)}
                  type="number"
                  className="form-control border border-black text-black px-2 py-1 flex-1 rounded"
                  id="price"
                  required
                />
              </div>
              <div className="mb-3 flex items-center gap-2">
                <label htmlFor="spots" className="text-black font-medium w-24">Max Spots:</label>
                <input
                  value={spots}
                  onChange={e => setSpots(e.target.value)}
                  type="number"
                  className="form-control border border-black text-black px-2 py-1 flex-1 rounded"
                  id="spots"
                  required
                />
              </div>
              <div className="flex justify-center mt-4">
                <button 
                  type="submit" 
                  className="bg-blue-600 text-white py-2 px-8 rounded hover:-translate-y-1 hover:scale-110 hover:bg-indigo-500 transition"
                >
                  {editingLot ? "Update" : "Add"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
    </>
  );
}

export default Admin;