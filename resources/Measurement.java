package resources;

import java.util.ArrayList;
import java.util.List;

public class Measurement {

    public static class TimePeriodType {
        public String startTime;
        public String endTime;
    }

    public static class MeasurementDataType {
        public String valueType;
        public String timestamp;
        public double value;
        public TimePeriodType evaluationPeriod;
        public String valueSource;
        public String valueTendency;
        public String valueState;
    }

    public static class MeasurementDescriptionDataType {
        public String measurementType;
        public String unit;
        public String scopeType;
        public String label;
        public String description;
    }

    public static List<FunctionModel> measurement(String measurementType, String unit, String scopeType, String label, String description) {
        List<FunctionModel> functionModels = new ArrayList<>();

        functionModels.add(new FunctionModel() {{
            functionName = "measurementData";
            function = new MeasurementDataType();
        }});

        functionModels.add(new FunctionModel() {{
            functionName = "measurementDescription";
            function = new MeasurementDescriptionDataType() {{
                this.measurementType = measurementType;
                this.unit = unit;
                this.scopeType = scopeType;
                this.label = label;
                this.description = description;
            }};
        }});

        return functionModels;
    }
}
