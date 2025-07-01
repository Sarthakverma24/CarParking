import { useEffect, useState } from "react";
import { Bar } from "react-chartjs-2";
import { useNavigate } from "react-router-dom";
import {
  Chart as ChartJS,
  BarElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend
} from "chart.js";

ChartJS.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend);

function SummeryUser() {
  const navigate = useNavigate();

  const [chartData, setChartData] = useState({
    labels: [],
    datasets: [
      {
        label: "Number of Past Bookings",
        data: [],
        backgroundColor: "#3b82f6",
        borderRadius: 5,
      },
    ],
  });

  useEffect(() => {
    fetch("http://localhost:8080/api/used-summary")
      .then((res) => res.json())
      .then((data) => {
        setChartData({
          labels: data.map((item) => item.location),
          datasets: [
            {
              label: "Number of Past Bookings",
              data: data.map((item) => item.count),
              backgroundColor: "#3b82f6",
              borderRadius: 5,
            },
          ],
        });
      })
      .catch((err) => {
        console.error("Failed to fetch used summary:", err);
      });
  }, []);

  return (
    <div className="min-h-screen bg-blue-950">
      <div className="flex items-center justify-between bg-blue-900 p-4 w-full">
        <h1 className="text-xl text-white font-semibold">Welcome Admin</h1>
        <div className="flex gap-4">
          <button
            className="px-4 py-2 bg-white rounded-lg shadow transition hover:bg-gray-100"
            onClick={() => navigate("/Dashboard")}
          >
            Home
          </button>
          <button
            className="px-4 py-2 bg-white rounded-lg shadow transition hover:bg-gray-100"
            onClick={() => navigate("/SummeryUser")}
          >
            Summary
          </button>
          <button
            className="px-4 py-2 bg-white rounded-lg shadow transition hover:bg-gray-100"
            onClick={() => navigate("/")}
          >
            Logout
          </button>
        </div>
        <h1 className="text-base underline cursor-pointer text-white hover:text-blue-300 transition">
          Edit Profile
        </h1>
      </div>

      <div className="flex justify-center items-center p-6">
        <div className="bg-white p-4 rounded-lg shadow w-full md:w-[600px]">
          <h2 className="text-center font-semibold mb-4 text-black">
            Used Parking Summary
          </h2>
          <Bar
            data={chartData}
            options={{
              responsive: true,
              plugins: {
                legend: {
                  position: "top",
                },
              },
              scales: {
                y: {
                  beginAtZero: true,
                  ticks: {
                    stepSize: 1,
                  },
                },
              },
            }}
          />
        </div>
      </div>
    </div>
  );
}

export default SummeryUser;
