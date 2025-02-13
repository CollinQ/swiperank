import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Users, Zap, BarChart3 } from "lucide-react"

export default function Features() {
    const features = [
        {
            icon: <Users className="h-8 w-8 text-blue500" />,
            title: "Intuitive Comparison",
            description: "Compare applicants head-to-head for fair and efficient evaluations."
        },
        {
            icon: <Zap className="h-8 w-8 text-blue-500" />,
            title: "Time-Saving",
            description: "Streamline your review process and save hours of manual work."
        },
        {
            icon: <BarChart3 className="h-8 w-8 text-blue-500" />,
            title: "Clear Insights",
            description: "Get a comprehensive ranking of all applicants based on your preferences."
        }
    ]

    return (
        <section id="features" className="py-20 px-6 md:px-10">
            <h2 className="text-3xl md:text-4xl font-bold text-center mb-12">Key Features</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                {features.map((feature, index) => (
                    <Card key={index}>
                        <CardHeader>
                            <CardTitle className="flex items-center">
                                {feature.icon}
                                <span className="m1-2">{feature.title}</span>
                            </CardTitle>
                        </CardHeader>
                        <CardContent>{feature.description}</CardContent>
                    </Card>
                ))}
            </div>
        </section>
    )
}