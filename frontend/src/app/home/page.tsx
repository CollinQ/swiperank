'use client';
import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { useRouter } from 'next/navigation';

// Mock data - would come from API
const mockProjects = [
    {
        id: 1,
        name: "Spring 2024 Club Applications",
        totalApplicants: 50,
        completedComparisons: 125,
        totalComparisons: 200,
        progress: 62.5
    },
    {
        id: 2, 
        name: "Fall 2023 Executive Board",
        totalApplicants: 20,
        completedComparisons: 80,
        totalComparisons: 100,
        progress: 80
    }
];

export default function Home() {
    const [projects, setProjects] = useState(mockProjects);
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        // API endpoint would be something like:
        // GET /api/projects
        // This would fetch active projects for the logged-in user
        const fetchProjects = async () => {
            try {
                // const response = await fetch('/api/projects');
                // const data = await response.json();
                // setProjects(data);
                setProjects(mockProjects); // Using mock data for now
                setLoading(false);
            } catch (error) {
                console.error('Error fetching projects:', error);
                setLoading(false);
            }
        };

        fetchProjects();
    }, []);

    const handleStartReview = (projectId: number) => {
        // This would navigate to the comparison interface for the selected project
        router.push(`/review/${projectId}`);
    };

    if (loading) {
        return <div className="flex justify-center items-center h-screen">Loading...</div>;
    }

    return (
        <div className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold mb-8">Active Projects</h1>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {projects.map((project) => (
                    <Card key={project.id} className="hover:shadow-lg transition-shadow">
                        <CardHeader>
                            <CardTitle>{project.name}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                <div>
                                    <p className="text-sm text-gray-500">Total Applicants: {project.totalApplicants}</p>
                                    <p className="text-sm text-gray-500">Comparisons: {project.completedComparisons} / {project.totalComparisons}</p>
                                </div>
                                <Progress value={project.progress} className="w-full" />
                                <Button 
                                    onClick={() => handleStartReview(project.id)}
                                    className="w-full"
                                >
                                    Continue Review
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}
