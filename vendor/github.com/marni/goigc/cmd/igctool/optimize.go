// Copyright ©2015 Ricardo Rocha <rocha.porto@gmail.com>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
	"github.com/marni/goigc"
	"log"

	"github.com/spf13/cobra"
)

var (
	method   string
	scorer   string
	tasktype []int
)

var optimizers = map[string]igc.Optimizer{
	"bf": igc.NewBruteForceOptimizer(false),
}

// optimizeCmd respresents the optimize command
var optimizeCmd = &cobra.Command{
	Use:   "optimize [path]",
	Short: "Optimize the track (for distance and score)",
	Long: `Parse a given track and return the maximum distance and score.
There is only one optimizer algorithm available at the moment:  
	bf: brute force
	
Available score functions:
  md1: max distance 1 turnpoint (out and return)
  md2: max distance with 2 turnpoints (triangle, non FAI)
  md3: max distance with 3 turnpoints
  md4: max distance with 4 turnpoints
  md5: max distance with 5 turnpoints
  fai: FAI defined triangle (shortest leg at least 28% of total triangle)
  nc: netcoupe.net (max 3 TPs, handicapped, -20% if not declared)
  olc: onlinecontest.org (max 5 TPs, handicapped, +30% if FAI triangle)
  
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("no track location given")
		}

		o, ok := optimizers[method]
		if !ok {
			log.Fatalf("unknown optimizer method: %v", method)
		}

		trk, err := igc.ParseLocation(args[0])
		if err != nil {
			log.Fatalf("optimization failed: %v", err)
		}

		r, err := o.Optimize(trk, len(trk.Points), igc.Distance)
		if err != nil {
			log.Fatalf("optimization failed: %v", err)
		}

		res := r.Distance()

		fmt.Printf("%#v\n", res)
	},
}

func init() {
	RootCmd.AddCommand(optimizeCmd)

	// Command flags
	optimizeCmd.PersistentFlags().StringVar(&scorer, "scorer", "sum", "score function (1TP, ..., 5TP, triangle, FAItriangle, netcoupe, olc)")
	optimizeCmd.PersistentFlags().StringVar(&method, "method", "bf", "optimizer method (bf, mc, ga)")
	optimizeCmd.PersistentFlags().IntSliceVarP(&tasktype, "tasktype", "t", []int{1}, "")

}
