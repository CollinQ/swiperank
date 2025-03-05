"use client";
import { useParams, useRouter } from "next/navigation";
import React, { useState, useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface FileInfo {
  fileID: string;
  fileName: string;
  fileType: string;
  data?: string;
}

interface Applicant {
  id: string;
  name: string;
  year: string;
  major: string;
  resume: FileInfo | null;
  coverLetter: FileInfo | null;
  image: FileInfo | null;
  elo: number;
  wins: number;
  losses: number;
}

const CandidatesPage = () => {
  const router = useRouter();
  const params = useParams();
  const projectId = params?.id as string;

  const [applicants, setApplicants] = useState<Applicant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchApplicants = async () => {
    try {
      console.log("Starting fetch...");
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      console.log(apiUrl);
      const response = await fetch(`${apiUrl}/api/getTwoForComparison`);

      console.log("Content-Type:", response.headers.get("content-type"));
      if (response.status === 409) {
        router.push(`/results/${projectId}`);
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to fetch applicants");
      }

      const data = await response.json();
      // Map _id(as stored in Mongo) to id
      const mappedData = data.map((a: any) => ({ ...a, id: a._id }));
      setApplicants(mappedData);
      setError(null);
    } catch (err: any) {
      console.error("Error fetching applicants:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApplicants();
  }, [projectId]);

  const handleCardSelect = async (winnerId: string, loserId: string) => {
    try {
      setLoading(true);
      const payload = {
        winnerId,
        loserId,
      };
      console.log("Sending payload:", payload); 
      
      const response = await fetch("http://localhost:8080/api/updateElo", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`Failed to update Elo ratings: ${await response.text()}`);
      }

      // Fetch next pair of applicants
      await fetchApplicants();
    } catch (err: any) {
      console.error("Error updating Elo:", err);
      setError(err.message);
      setLoading(false);
    }
  };

  const handleFileClick = (
    fileInfo: FileInfo | null,
    preview: boolean = false
  ) => {
    if (!fileInfo?.data) return;

    if (preview) {
      try {
        // Create data URL directly for PDF preview
        const dataUrl = `data:application/pdf;base64,${fileInfo.data}`;

        // Open PDF in new window/tab
        const newWindow = window.open("", "_blank");
        if (newWindow) {
          newWindow.document.write(`
            <html>
              <head>
                <title>${fileInfo.fileName}</title>
              </head>
              <body style="margin:0;padding:0;">
                <embed 
                  width="100%" 
                  height="100%" 
                  src="${dataUrl}" 
                  type="application/pdf"
                />
              </body>
            </html>
          `);
        }
      } catch (error) {
        console.error("Error processing PDF:", error);
        alert("Error opening PDF. Please try downloading instead.");
      }
    } else {
      const linkElement = document.createElement("a");
      linkElement.href = `data:${fileInfo.fileType};base64,${fileInfo.data}`;
      linkElement.download = fileInfo.fileName;
      document.body.appendChild(linkElement);
      linkElement.click();
      document.body.removeChild(linkElement);
    }
  };
  console.log(applicants);
  const [selectedId, setSelectedId] = useState();
  if (!applicants.length) {
    return <div>Loading...</div>;
  }

  return (
    <div className="flex flex-col items-center justify-center p-4 max-w-7xl mx-auto">
      {loading && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded-lg shadow-lg">
            <p className="text-lg font-medium">Loading next comparison...</p>
          </div>
        </div>
      )}
      <div className="flex flex-col md:flex-row gap-4">
        {applicants.map((applicant, index) => {
          return (
            <div key={`applicant-${applicant.id}-${index}`}>
              <Card
                onClick={() => {
                  console.log("Winner ID:", applicant.id);
                  console.log("Loser ID:", applicants[1 - index].id);
                  handleCardSelect(
                    applicant.id,
                    applicants[1 - index].id
                  )
                }}
                className={`w-full md:w-96 cursor-pointer transition-shadow border-2 hover:border-green-500 ${
                  selectedId === applicant.id
                    ? "shadow-xl border-blue-500"
                    : "shadow-sm border-gray-200"
                }`}
              >
                <CardHeader>
                  <CardTitle className="text-xl font-bold">
                    {applicant.name}
                  </CardTitle>
                  <div className="text-sm text-gray-600">
                    {applicant.year} • {applicant.major}
                  </div>
                  <div className="text-sm text-gray-600 mt-1">
                    Elo: {applicant.elo || "N/A"} ({applicant.wins || 0}W-{applicant.losses || 0}L)
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {applicant.image && (
                      <div>
                        <img
                          src={`data:${applicant.image.fileType};base64,${applicant.image.data}`}
                          alt={applicant.name}
                          className="w-full h-48 object-cover rounded-lg"
                        />
                      </div>
                    )}
                    <div>
                      <p className="font-semibold">Resume</p>
                      {applicant.resume ? (
                        <div className="space-x-2">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleFileClick(applicant.resume, true);
                            }}
                            className="text-blue-500 hover:underline"
                          >
                            View Resume
                          </button>
                          <span>•</span>
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleFileClick(applicant.resume);
                            }}
                            className="text-blue-500 hover:underline"
                          >
                            Download
                          </button>
                        </div>
                      ) : (
                        <p className="text-gray-600">No resume available</p>
                      )}
                    </div>
                    <Separator className="my-2" />
                    <div>
                      <p className="font-semibold">Cover Letter</p>
                      {applicant.coverLetter ? (
                        <div className="space-x-2">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleFileClick(applicant.coverLetter, true);
                            }}
                            className="text-blue-500 hover:underline"
                          >
                            View Cover Letter
                          </button>
                          <span>•</span>
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleFileClick(applicant.coverLetter);
                            }}
                            className="text-blue-500 hover:underline"
                          >
                            Download
                          </button>
                        </div>
                      ) : (
                        <p className="text-gray-600">No cover letter available</p>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default CandidatesPage;
