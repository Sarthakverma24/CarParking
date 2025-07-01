import { useState, useEffect } from "react";
import {useLocation, useNavigate } from "react-router-dom";
import Swal from "sweetalert2";

function Dashboard() {
  const [users, setUsers] = useState([]);
  const [slot, setslot] = useState([]);
  const [input, setinput] = useState("");
  const [vehicle_no,setvehicleno]=useState("");
  const [searchType, setSearchType] = useState("");
  const [popup, setPopup] = useState(false);
  const [selectedSlot, setSelectedSlot] = useState(null);
  const locations = useLocation();
  const { username } = locations.state || {};

  const navigate = useNavigate();

  useEffect(() => {
    fetchHistory();
  }, []);

  async function fetchHistory() {
    try {
      const response = await fetch(`http://localhost:8080/api/UserHistory/${username}`);
      if (!response.ok) throw new Error("Failed to fetch user history");
      const data = await response.json();
      setUsers(data ?? []);
      console.log(data);
      
    } catch (error) {
      Swal.fire({ icon: "error", title: "Error", text: error.message });
    }
  }

  async function getSlots() {
    if (!searchType || !input) {
      Swal.fire({
        icon: "warning",
        title: "Missing Input",
        text: "Please select a search type and enter a value",
      });
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/searchSlots?type=${searchType}&value=${input}`);
      console.log(searchType);
      
      if (!response.ok) throw new Error("Failed to fetch slots");
      const data = await response.json();
      setslot(data ?? []);
    } catch (error) {
      Swal.fire({ icon: "error", title: "Error", text: error.message });
    }
  }

  async function bookSlot(bookingData) {
    try {
      const response = await fetch(`http://localhost:8080/api/bookSlot`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(bookingData),
      });

      if (!response.ok) throw new Error("Failed to book slot");

      Swal.fire({ icon: "success", title: "Success", text: "Slot booked successfully!" });
      getSlots();
      fetchHistory();  
    } catch (error) {
      Swal.fire({ icon: "error", title: "Error", text: error.message });
    }
  }
  async function handleRelease(reservationData) {
    try {
      const response = await fetch(`http://localhost:8080/api/releaseSlot`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(reservationData),
      });

      if (!response.ok) throw new Error("Failed to release the slot");

      const updatedBooking = await response.json(); 
      console.log(updatedBooking);
      
      Swal.fire({ 
        icon: "success", 
        title: "Slot Released", 
        html: `
          <p><strong>Vehicle:</strong> ${updatedBooking.vehicle_no}</p>
          <p><strong>Slot ID:</strong> ${updatedBooking.spot_id}</p>
          <p><strong>Parking Cost:</strong> ₹${updatedBooking.parking_cost}</p>
        `
      });

      fetchHistory();  // Refresh the table
    } catch (err) {
      Swal.fire({ icon: "error", title: "Error", text: err.message });
    }
  }




  return (
    <>
      <div className="min-h-screen bg-blue-950">
        {/* Top Navbar */}
        <div className="flex items-center justify-between bg-blue-900 p-4 w-full">
          <h1 className="text-xl text-white font-semibold">Welcome User</h1>
          <div className="flex gap-4">
            <button className="px-4 py-2 bg-white rounded-lg shadow" onClick={() => navigate("/Dashboard")}>Home</button>
            <button className="px-4 py-2 bg-white rounded-lg shadow"
            onClick={() =>navigate("/SummeryUser")}>Summary</button>
            <button className="px-4 py-2 bg-white rounded-lg shadow" onClick={() => navigate("/")}>Logout</button>
          </div>
          <h1 className="text-base underline cursor-pointer hover:text-blue-700">Edit Profile</h1>
        </div>

        {/* Parking History */}
        <div className="p-4 bg-black flex justify-center border text-white border-black m-5 rounded-lg">
          <h5 className="text-lg font-bold">Recent Parking History</h5>
          </div>

          <div className="p-4 bg-blue-900 border border-black m-5 rounded-lg" style={{ height: "400px" }}>
            <div className="overflow-y-auto h-full">
              <table className="min-w-full text-white text-left border-collapse">
                <thead>
                  <tr className="bg-blue-800 text-white">
                    <th className="px-4 py-2 border-b border-gray-600">ID</th>
                    <th className="px-4 py-2 border-b border-gray-600">Slot ID</th>
                    <th className="px-4 py-2 border-b border-gray-600">Location</th>
                    <th className="px-4 py-2 border-b border-gray-600">Vehicle Number</th>
                    <th className="px-4 py-2 border-b border-gray-600">Parking Time</th>
                    <th className="px-4 py-2 border-b border-gray-600">Leaving Time</th>
                    <th className="px-4 py-2 border-b border-gray-600">Status</th>
                  </tr>
                </thead>
                <tbody>
                {users.length === 0 ? (
                  <tr><td colSpan="7" className="text-center py-4">No history found.</td></tr>
                ) : (
                  [...users]
                    .sort((a, b) => (b.status === true) - (a.status === true)) // true rows first
                    .map((user, index) => (
                      <tr key={index} className="hover:bg-blue-700 transition">
                        <td className="px-4 py-2 border-b border-gray-700">{user.id}</td>
                        <td className="px-4 py-2 border-b border-gray-700">{user.spot_id}</td>
                        <td className="px-4 py-2 border-b border-gray-700">{user.location}</td>
                        <td className="px-4 py-2 border-b border-gray-700">{user.vehicle_no}</td>
                        <td className="px-4 py-2 border-b border-gray-700">{new Date(user.parking).toLocaleString()}</td>
                        <td className="px-4 py-2 border-b border-gray-700">
                          {user.status ? "Occupied" : new Date(user.leaving).toLocaleString()}
                        </td>
                        <td className="px-4 py-2 border-b border-gray-700">
                          {user.status ? (
                            <button
                              className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                              onClick={() => handleRelease({ id: user.id, spot_id: user.spot_id })}
                            >
                              Release
                            </button>
                          ) : (
                            "Released"
                          )}
                        </td>
                      </tr>
                    ))
                )}
              </tbody>
              </table>
            </div>
          </div>
 
        {/* Search Slots */}
       <div className="w-full px-4 mt-6">
        <label htmlFor="text" className="block text-white mb-2 font-semibold">Search with</label>
        <div className="flex flex-col sm:flex-row gap-4">
          {/* Dropdown first */}
          <select
            onChange={e => setSearchType(e.target.value)}
            id="search"
            className="p-2 rounded-md bg-white text-black border border-gray-300 sm:w-1/2 w-full focus:outline-none"
          >
            <option value="">Choose an option</option>
            <option value="location">Location</option>
            <option value="pincode">Pincode</option>
          </select>

          {/* Input second */}
          <input
            onChange={e => setinput(e.target.value)}
            onKeyDown={e => {
              if (e.key === "Enter") {
                e.preventDefault();
                getSlots();
              }
            }}
            type="text"
            id="text"
            className="bg-blue-200 text-black px-3 py-2 rounded-md border border-gray-300 sm:w-1/2 w-full"
            placeholder="Enter details"
          />
        </div>

        {/* Search Button */}
        <div className="mt-4">
          <button
            onClick={getSlots}
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 py-2 rounded-md"
          >
            Search
          </button>
        </div>
      </div>

        {/* Slot List */}
        <div className="p-4 bg-blue-900 border border-black m-5 rounded-lg overflow-x-auto">
          <table className="min-w-full text-white text-left border-collapse">
            <thead>
              <tr className="bg-blue-800 text-white">
                <th className="px-4 py-2 border-b border-gray-600">ID</th>
                <th className="px-4 py-2 border-b border-gray-600">Location</th>
                <th className="px-4 py-2 border-b border-gray-600">Availability</th>
                <th className="px-4 py-2 border-b border-gray-600">Action</th>
              </tr>
            </thead>
            <tbody>
              {slot.length === 0 ? (
                <tr><td colSpan="4" className="text-center py-4">Search Parking Slots</td></tr>
              ) : (
                slot.map((s, index) => (
                  <tr key={index} className="hover:bg-blue-700 transition">
                    <td className="px-4 py-2 border-b border-gray-700">{s.id}</td>
                    <td className="px-4 py-2 border-b border-gray-700">{s.location}</td>
                    <td className="px-4 py-2 border-b border-gray-700">{s.spot_list?.filter(spot => spot.status === true).length ?? 0}</td>
                    <td className="px-4 py-2 border-b border-gray-700">
                      <button
                        onClick={() => {
                          const available = s.spot_list?.find(spot => spot.status === true);
                            if (available) {
                              setSelectedSlot({ ...available, location: s.location });
                              setPopup(true);
                            }  else {
                            Swal.fire({
                              icon: "info",
                              title: "No Available Spot",
                              text: "All spots in this lot are occupied.",
                            });
                          }
                        }}
                        disabled={!s.spot_list?.some(spot => spot.status === true)}
                        className={`px-3 py-1 rounded ${
                          s.spot_list?.some(spot => spot.status === true)
                            ? "bg-green-500 hover:bg-green-600 text-white"
                            : "bg-gray-400 cursor-not-allowed text-white"
                        }`}
                      >
                        Book
                      </button> 
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
        {popup && selectedSlot && (
  <div className="fixed inset-0 flex items-center justify-center z-50">
    <div className="bg-white rounded-2xl p-6 w-[400px] max-w-full shadow-xl relative border border-gray-300">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-bold text-gray-800">Book the Parking Spot</h2>
        <button
          className="text-gray-500 hover:text-gray-800 font-bold text-xl"
          onClick={() => {
            setPopup(false);
            setSelectedSlot(null);
          }}
        >
          ×
        </button>
      </div>
          <form
              onSubmit={(e) => {
                e.preventDefault();

                const bookingData = {
                  lot_id: selectedSlot.lot_id,
                  spot_id: selectedSlot.parking_id,
                  user_id: username,
                  vehicle_no: vehicle_no,
                  location: selectedSlot.location,
                };

                bookSlot(bookingData);
                setPopup(false);
                setSelectedSlot(null);
              }}
              className="w-full space-y-4"
            >
            <div className="flex items-center gap-2">
              <label htmlFor="slotId" className="text-black font-medium w-24">Slot ID:</label>
              <input
                type="text"
                id="slotId"
                value={selectedSlot.parking_id}
                readOnly
                className="form-control border border-black bg-gray-100 text-black px-2 py-1 flex-1 rounded"
              />
            </div>
            <div className="flex items-center gap-2">
              <label htmlFor="lotId" className="text-black font-medium w-24">lot ID:</label>
              <input
                type="text"
                id="lotId"
                value={selectedSlot.lot_id}
                readOnly
                className="form-control border border-black bg-gray-100 text-black px-2 py-1 flex-1 rounded"
              />
            </div>
            <div className="flex items-center gap-2">
              <label htmlFor="userId" className="text-black font-medium w-24">User ID:</label>
              <input
                type="text"
                id="userId"
                value={username}
                readOnly
                className="form-control border border-black bg-gray-100 text-black px-2 py-1 flex-1 rounded"
              />
            </div>
            <div className="flex items-center gap-2">
              <label htmlFor="status" className="text-black font-medium w-24">Status:</label>
              <input
                type="text"
                id="status"
                value={selectedSlot.status ? "Vacant" : "Occupied"}
                readOnly
                className="form-control border border-black bg-gray-100 text-black px-2 py-1 flex-1 rounded"
              />
            </div>
            <div className="flex items-center gap-2">
              <label htmlFor="userId" className="text-black font-medium w-24">Vehicle Number:</label>
              <input
                type="text"
                id="userId"
                onChange={e => {setvehicleno(e.target.value)}}
                className="form-control border border-black bg-gray-100 text-black px-2 py-1 flex-1 rounded"
              />
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
              >
                Confirm Booking
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

export default Dashboard;
