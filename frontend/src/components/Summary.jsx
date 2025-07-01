import { Doughnut, Pie } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";
import Swal from "sweetalert2";

ChartJS.register(ArcElement, Tooltip, Legend, ChartDataLabels);

function Summary() {
  const [revenueData, setRevenueData] = useState({
    labels: [],
    datasets: [{ label: "Revenue", data: [], backgroundColor: [] }]
  });
  const [lotStatusData, setLotStatusData] = useState({
    labels: [],
    datasets: [{ label: "Status", data: [], backgroundColor: [] }]
  });

  const navigate = useNavigate();

  useEffect(() => {
    fetch("http://localhost:8080/api/summary")
      .then(res => res.json())
      .then(data => {
        setRevenueData({
          labels: data.revenue.map(lot => lot.location),
          datasets: [{
            label: "Revenue",
            data: data.revenue.map(lot => lot.amount),
            backgroundColor: ['#22c55e', '#3b82f6', '#f59e0b', '#ef4444'],
          }]
        });

        setLotStatusData({
          labels: ["Occupied", "Available"],
          datasets: [{
            label: "Status",
            data: [data.occupied, data.available],
            backgroundColor: ['#ef4444', '#22c55e'],
          }]
        });
      })
      .catch(err => {
        console.error("Failed to fetch summary", err);
        Swal.fire({ icon: "error", title: "Error", text: err.message });
      });
  }, []);

  const optionsWithLabels = {
    plugins: {
      datalabels: {
        color: '#000',
        font: { weight: 'bold' },
        formatter: (value) => value
      },
      legend: {
        position: 'bottom',
         labels: {
        color: '#000', 
        font: {
          size: 15
        }
        }
      }
    }
  };

  return (
    <div className="min-h-screen bg-blue-950">
      <div className="flex items-center justify-between bg-blue-900 p-4 w-full">
        <h1 className="text-xl text-white font-semibold">Welcome Admin</h1>
        <div className="flex gap-4">
          <button className="px-4 py-2 bg-white rounded-lg shadow transition" onClick={() => navigate('/Admin')}>Home</button>
          <button className="px-4 py-2 bg-white rounded-lg shadow transition" onClick={() => navigate('/Users')}>Users</button>
          <button className="px-4 py-2 bg-white rounded-lg shadow transition" onClick={() => navigate('/Summary')}>Summary</button>
          <button className="px-4 py-2 bg-white rounded-lg shadow transition" onClick={() => navigate('/')}>Logout</button>
        </div>
        <h1 className="text-base underline cursor-pointer hover:text-blue-700 transition-colors">Edit Profile</h1>
      </div>

      <div className="flex flex-wrap justify-center items-center gap-8 p-15">
        {revenueData?.datasets?.[0]?.data?.length > 0 && (
          <div className="bg-white p-4 rounded-lg shadow w-[350px] h-[350px]">
            <h2 className="text-center font-semibold mb-2 text-black">Revenue by Lot</h2>
            <Doughnut data={revenueData} options={optionsWithLabels} plugins={[ChartDataLabels]} />
          </div>
        )}

        {lotStatusData?.datasets?.[0]?.data?.length > 0 && (
          <div className="bg-white p-4 rounded-lg shadow w-[350px] h-[350px]">
            <h2 className="text-center font-semibold mb-2 text-black">Slot Occupancy</h2>
            <Pie data={lotStatusData} options={optionsWithLabels} plugins={[ChartDataLabels]} />
          </div>
        )}
      </div>
    </div>
  );
}

export default Summary;
