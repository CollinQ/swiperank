"use client";

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
  ratingCount: number;
  rating: number;
}

const CandidatesPage = () => {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [applicants, setApplicants] = useState<Applicant[]>([]);

  useEffect(() => {
    const fetchApplicants = async () => {
      try {
        console.log('Starting fetch...');
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
        const response = await fetch(`${apiUrl}/least-rated-applicants`, {
          method: 'GET',
          mode: 'cors',
          headers: {
            'Accept': 'application/json'
          }
        });
        console.log('Status:', response.status);
        console.log('Content-Type:', response.headers.get('content-type'));
        
        if (!response.ok) {
          const text = await response.text();
          console.error('Error response:', text.substring(0, 200) + '...');
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        
        setApplicants(result.data.applicants || []);
      } catch (error) {
        console.error('Fetch error details:', error);
        setApplicants([]);
      }
    };
    fetchApplicants();
  }, []);

  const handleFileClick = (fileInfo: FileInfo | null, preview: boolean = false) => {
    if (!fileInfo?.data) return;
    
    
    if (preview) {
      try {
        // Create data URL directly for PDF preview
        const dataUrl = `data:application/pdf;base64,${fileInfo.data}`;
        
        // Open PDF in new window/tab
        const newWindow = window.open('', '_blank');
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
        console.error('Error processing PDF:', error);
        alert('Error opening PDF. Please try downloading instead.');
      }
    } else {
      // Original download functionality
      const linkElement = document.createElement('a');
      linkElement.href = `data:${fileInfo.fileType};base64,${fileInfo.data}`;
      linkElement.download = fileInfo.fileName;
      document.body.appendChild(linkElement);
      linkElement.click();
      document.body.removeChild(linkElement);
    }
  };

  const [leftApplicant, rightApplicant] = applicants;

  if (!leftApplicant || !rightApplicant) {
    return <div>Loading...</div>;
  }

  return (
    <div className="flex flex-col items-center justify-center p-4 max-w-7xl mx-auto">
      <div className="flex flex-col md:flex-row gap-4">
        {[leftApplicant, rightApplicant].map((applicant, index) => (
          <div key={`applicant-${applicant.id}-${index}`}>
            <Card
              onClick={() => setSelectedId(applicant.id)}
              className={`w-full md:w-96 cursor-pointer transition-shadow border-2 ${
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
                  Rating: {applicant.rating?.toFixed(1) || 'N/A'} ({applicant.ratingCount} reviews)
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
        ))}
      </div>
    </div>
  );
};

export default CandidatesPage;
