"use client";

import { useParams } from "next/navigation";
import { useState, useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Crown, Medal } from "lucide-react";
import Link from "next/link";

interface ApplicantRanking {
  id: string;
  first_name: string;
  last_name: string;
  elo: number;
  wins: number;
  losses: number;
}

export default function ResultsPage() {
  const params = useParams();
  const projectId = params?.id as string;

  const [rankings, setRankings] = useState<ApplicantRanking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRankings = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/rankings?project_id=${projectId}`
        );
        if (!response.ok) {
          throw new Error("Failed to fetch rankings");
        }
        const data = await response.json();
        setRankings(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchRankings();
  }, [projectId]);

  if (loading) return <div className="text-center">Loading rankings...</div>;
  if (error) return <div className="text-center text-red-500">{error}</div>;

  const topThree = rankings.slice(0, 3);
  const restOfRankings = rankings.slice(3);

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-12 text-center">Final Rankings</h1>

      {/* Podium Section */}
      <div className="flex justify-center items-end gap-4 mb-16 h-[400px]">
        {/* Second Place */}
        {topThree[1] && (
          <div className="flex flex-col items-center w-64">
            <Card className="w-full mb-2 bg-gradient-to-b from-gray-100 to-gray-300">
              <CardHeader>
                <Medal className="h-8 w-8 text-gray-500 mx-auto" />
                <CardTitle className="text-center">
                  {topThree[1].first_name} {topThree[1].last_name}
                </CardTitle>
              </CardHeader>
              <CardContent className="text-center">
                <p className="text-gray-600">Elo: {topThree[1].elo}</p>
                <p className="text-sm">
                  W: {topThree[1].wins} - L: {topThree[1].losses}
                </p>
              </CardContent>
            </Card>
            <div className="h-32 w-full bg-gray-300 rounded-t-lg" />
          </div>
        )}

        {/* First Place */}
        {topThree[0] && (
          <div className="flex flex-col items-center w-64">
            <Crown className="h-12 w-12 text-yellow-500 mb-2" />
            <Link href={`/applicants/${topThree[0].id}`}>
              <Card className="w-full mb-2 bg-gradient-to-b from-yellow-100 to-yellow-300 hover:shadow-lg transition-shadow">
                <CardHeader>
                  <Medal className="h-8 w-8 text-yellow-500 mx-auto" />
                  <CardTitle className="text-center">
                    {topThree[0].first_name} {topThree[0].last_name}
                  </CardTitle>
                </CardHeader>
                <CardContent className="text-center">
                  <p className="text-gray-600">Elo: {topThree[0].elo}</p>
                  <p className="text-sm">
                    W: {topThree[0].wins} - L: {topThree[0].losses}
                  </p>
                </CardContent>
              </Card>
            </Link>
            <div className="h-40 w-full bg-yellow-300 rounded-t-lg" />
          </div>
        )}

        {/* Third Place */}
        {topThree[2] && (
          <div className="flex flex-col items-center w-64">
            <Card className="w-full mb-2 bg-gradient-to-b from-orange-100 to-orange-300">
              <CardHeader>
                <Medal className="h-8 w-8 text-orange-500 mx-auto" />
                <CardTitle className="text-center">
                  {topThree[2].first_name} {topThree[2].last_name}
                </CardTitle>
              </CardHeader>
              <CardContent className="text-center">
                <p className="text-gray-600">Elo: {topThree[2].elo}</p>
                <p className="text-sm">
                  W: {topThree[2].wins} - L: {topThree[2].losses}
                </p>
              </CardContent>
            </Card>
            <div className="h-24 w-full bg-orange-300 rounded-t-lg" />
          </div>
        )}
      </div>

      {/* Rest of Rankings */}
      <div className="grid gap-4 max-w-2xl mx-auto">
        {restOfRankings.map((applicant, index) => (
          <Link key={applicant.id} href={`/applicants/${applicant.id}`}>
            <Card className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <CardTitle className="flex justify-between items-center">
                  <span>
                    #{index + 4} - {applicant.first_name} {applicant.last_name}
                  </span>
                  <span className="text-sm text-gray-500">
                    Elo: {applicant.elo}
                  </span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex justify-between text-sm text-gray-600">
                  <span>Wins: {applicant.wins}</span>
                  <span>Losses: {applicant.losses}</span>
                </div>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}