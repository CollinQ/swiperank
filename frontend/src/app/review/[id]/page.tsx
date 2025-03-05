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
  _id: string;
  first_name: string;
  last_name: string;
  year: string;
  major: string;
  resume: FileInfo | null;
  coverLetter: FileInfo | null;
  image: FileInfo | null;
  ratingCount: number;
  rating: number;
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
      const response = await fetch(`${apiUrl}/api/getTwoForComparison`, {
        method: "GET",
        mode: "cors",
        headers: {
          Accept: "application/json",
        },
      });

      if (response.status === 409) {
        router.push(`/results/${projectId}`);
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to fetch applicants");
      }

      const data: Applicant[] = await response.json();

      setApplicants(data);
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
    console.log(winnerId);
    console.log(loserId);
    try {
      // Update Elo ratings
      const response = await fetch("http://localhost:8080/api/updateElo", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          winner_id: winnerId,
          loser_id: loserId,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to update Elo ratings");
      }

      // Fetch next pair of applicants
      await fetchApplicants();
    } catch (err: any) {
      console.error("Error updating Elo:", err);
      setError(err.message);
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
      // Original download functionality
      const linkElement = document.createElement("a");
      linkElement.href = `data:${fileInfo.fileType};base64,${fileInfo.data}`;
      linkElement.download = fileInfo.fileName;
      document.body.appendChild(linkElement);
      linkElement.click();
      document.body.removeChild(linkElement);
    }
  };
  console.log(applicants);
  const [leftApplicant, rightApplicant] = applicants;
  const [selectedId, setSelectedId] = useState();
  if (!leftApplicant || !rightApplicant) {
    return <div>Loading...</div>;
  }

  return (
    <div className="flex flex-col items-center justify-center p-4 max-w-7xl mx-auto">
      <div className="flex flex-col md:flex-row gap-4">
        <div>
          <Card
            onClick={() => {
              handleCardSelect(leftApplicant._id, rightApplicant._id);
            }}
            className={`w-full md:w-96 cursor-pointer transition-shadow border-2 ${
              selectedId === leftApplicant._id
                ? "shadow-xl border-blue-500"
                : "shadow-sm border-gray-200"
            }`}
          >
            <CardHeader>
              <CardTitle className="text-xl font-bold">
                {leftApplicant.first_name} {leftApplicant.last_name}
              </CardTitle>
              <div className="text-sm text-gray-600">
                {leftApplicant.year} • {leftApplicant.major}
              </div>
              <div className="text-sm text-gray-600 mt-1">
                Rating: {leftApplicant.rating?.toFixed(1) || "N/A"} (
                {leftApplicant.ratingCount} reviews)
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {leftApplicant.image && (
                  <div>
                    <img
                      src={`data:${leftApplicant.image.fileType};base64,${leftApplicant.image.data}`}
                      alt={leftApplicant.name}
                      className="w-full h-48 object-cover rounded-lg"
                    />
                  </div>
                )}
                <div>
                  <p className="font-semibold">Resume</p>
                  {leftApplicant.resume ? (
                    <div className="space-x-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(leftApplicant.resume, true);
                        }}
                        className="text-blue-500 hover:underline"
                      >
                        View Resume
                      </button>
                      <span>•</span>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(leftApplicant.resume);
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
                  {leftApplicant.coverLetter ? (
                    <div className="space-x-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(leftApplicant.coverLetter, true);
                        }}
                        className="text-blue-500 hover:underline"
                      >
                        View Cover Letter
                      </button>
                      <span>•</span>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(leftApplicant.coverLetter);
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

        <div>
          <Card
            onClick={() =>
              handleCardSelect(rightApplicant._id, leftApplicant._id)
            }
            className={`w-full md:w-96 cursor-pointer transition-shadow border-2 ${
              selectedId === rightApplicant._id
                ? "shadow-xl border-blue-500"
                : "shadow-sm border-gray-200"
            }`}
          >
            <CardHeader>
              <CardTitle className="text-xl font-bold">
                {rightApplicant.first_name} {rightApplicant.last_name}
              </CardTitle>
              <div className="text-sm text-gray-600">
                {rightApplicant.year} • {rightApplicant.major}
              </div>
              <div className="text-sm text-gray-600 mt-1">
                Rating: {rightApplicant.rating?.toFixed(1) || "N/A"} (
                {rightApplicant.ratingCount} reviews)
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {rightApplicant.image && (
                  <div>
                    <img
                      src={`data:${rightApplicant.image.fileType};base64,${rightApplicant.image.data}`}
                      alt={rightApplicant.first_name}
                      className="w-full h-48 object-cover rounded-lg"
                    />
                  </div>
                )}
                <div>
                  <p className="font-semibold">Resume</p>
                  {rightApplicant.resume ? (
                    <div className="space-x-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(rightApplicant.resume, true);
                        }}
                        className="text-blue-500 hover:underline"
                      >
                        View Resume
                      </button>
                      <span>•</span>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(rightApplicant.resume);
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
                  {rightApplicant.coverLetter ? (
                    <div className="space-x-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(rightApplicant.coverLetter, true);
                        }}
                        className="text-blue-500 hover:underline"
                      >
                        View Cover Letter
                      </button>
                      <span>•</span>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleFileClick(rightApplicant.coverLetter);
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
      </div>
    </div>
  );
};

export default CandidatesPage;
